package server

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"mcp-middleware/config"
)

// aliasClaim is the JWT payload claim that identifies the tenant. The tenant data-plane
// base URL is derived from it via cfg.TenantBaseURLTemplate (default
// "https://{alias}.middleware.io"). Fixed by contract with app.middleware.io.
const aliasClaim = "alias"

// aliasPattern restricts an alias to a single safe DNS label so it can't inject a host,
// path, scheme, or port when substituted into the URL template.
var aliasPattern = regexp.MustCompile(`^[a-zA-Z0-9]([a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?$`)

type ctxKey string

const (
	ctxKeyToken         ctxKey = "mw_bearer_token"
	ctxKeyTenantBaseURL ctxKey = "mw_tenant_base_url_ctx"
	ctxKeyExcludedTools ctxKey = "mw_excluded_tools"
)

func WithToken(ctx context.Context, token string) context.Context {
	return context.WithValue(ctx, ctxKeyToken, token)
}

func TokenFromContext(ctx context.Context) (string, bool) {
	v, ok := ctx.Value(ctxKeyToken).(string)
	return v, ok && v != ""
}

func WithTenantBaseURL(ctx context.Context, u string) context.Context {
	return context.WithValue(ctx, ctxKeyTenantBaseURL, u)
}

func TenantBaseURLFromContext(ctx context.Context) (string, bool) {
	v, ok := ctx.Value(ctxKeyTenantBaseURL).(string)
	return v, ok && v != ""
}

func WithExcludedTools(ctx context.Context, set map[string]struct{}) context.Context {
	return context.WithValue(ctx, ctxKeyExcludedTools, set)
}

func ExcludedToolsFromContext(ctx context.Context) map[string]struct{} {
	v, _ := ctx.Value(ctxKeyExcludedTools).(map[string]struct{})
	return v
}

func parseExcludedTools(rawQuery string) map[string]struct{} {
	values, err := url.ParseQuery(rawQuery)
	if err != nil {
		return nil
	}
	raw := values.Get("exclude_tools")
	if raw == "" {
		return nil
	}
	out := make(map[string]struct{})
	for _, name := range strings.Split(raw, ",") {
		name = strings.TrimSpace(name)
		if name != "" {
			out[name] = struct{}{}
		}
	}
	return out
}

// ProtectedResourceMetadataHandler serves the RFC 9728 OAuth 2.0 Protected Resource Metadata
// document for this MCP server. Claude fetches this after a 401 to discover the AS.
func ProtectedResourceMetadataHandler(cfg *config.Config) http.HandlerFunc {
	doc := map[string]any{
		"resource":                 cfg.MCPServerURL,
		"authorization_servers":    []string{cfg.AuthServerURL},
		"scopes_supported":         strings.Fields(cfg.MCPScopes),
		"bearer_methods_supported": []string{"header"},
		"resource_name":            "Middleware MCP Server",
	}
	body, _ := json.Marshal(doc)
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Cache-Control", "public, max-age=300")
		_, _ = w.Write(body)
	}
}

// RequireBearer enforces the presence of an Authorization: Bearer token, decodes the JWT payload
// (no signature verification — backend remains the authority), and extracts the tenant base URL
// from the configured claim. On failure it returns 401 with a WWW-Authenticate header pointing at
// the Protected Resource Metadata document, per the MCP authorization spec.
func RequireBearer(cfg *config.Config, next http.Handler) http.Handler {
	resourceMetadataURL := buildResourceMetadataURL(cfg)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := extractBearerToken(r.Header.Get("Authorization"))
		if token == "" {
			writeUnauthorized(w, resourceMetadataURL, cfg.MCPScopes, "", "")
			return
		}

		tenantURL, err := tenantBaseURLFromJWT(token, cfg.TenantBaseURLTemplate)
		if err != nil {
			writeUnauthorized(w, resourceMetadataURL, cfg.MCPScopes, "invalid_token", err.Error())
			return
		}

		ctx := r.Context()
		ctx = WithToken(ctx, token)
		ctx = WithTenantBaseURL(ctx, tenantURL)
		if excluded := parseExcludedTools(r.URL.RawQuery); len(excluded) > 0 {
			ctx = WithExcludedTools(ctx, excluded)
		}
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func extractBearerToken(authHeader string) string {
	if authHeader == "" {
		return ""
	}
	const prefix = "bearer "
	if len(authHeader) < len(prefix) || !strings.EqualFold(authHeader[:len(prefix)], prefix) {
		return ""
	}
	return strings.TrimSpace(authHeader[len(prefix):])
}

// tenantBaseURLFromJWT decodes the (unverified) JWT payload, reads the tenant `alias`
// claim, and builds the tenant data-plane base URL by substituting it into urlTemplate
// (e.g. "https://{alias}.middleware.io" → "https://acme.middleware.io"). The alias is
// validated as a single DNS label first so it can't inject a host/path/scheme.
func tenantBaseURLFromJWT(token, urlTemplate string) (string, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return "", fmt.Errorf("access token is not a JWT")
	}
	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		// Some issuers pad with '='; try standard URL encoding too.
		payload, err = base64.URLEncoding.DecodeString(parts[1])
		if err != nil {
			return "", fmt.Errorf("failed to base64-decode JWT payload")
		}
	}
	var claims map[string]any
	if err := json.Unmarshal(payload, &claims); err != nil {
		return "", fmt.Errorf("failed to parse JWT payload as JSON")
	}
	raw, ok := claims[aliasClaim]
	if !ok {
		return "", fmt.Errorf("JWT missing required claim %q", aliasClaim)
	}
	alias, ok := raw.(string)
	if !ok || alias == "" {
		return "", fmt.Errorf("JWT claim %q is not a non-empty string", aliasClaim)
	}
	if !aliasPattern.MatchString(alias) {
		return "", fmt.Errorf("JWT claim %q is not a valid tenant alias", aliasClaim)
	}

	base := strings.ReplaceAll(urlTemplate, "{alias}", alias)
	u, err := url.Parse(base)
	if err != nil || u.Host == "" {
		return "", fmt.Errorf("derived tenant base URL %q is not valid", base)
	}
	if u.Scheme != "https" && !(u.Scheme == "http" && isLoopbackHost(u.Hostname())) {
		return "", fmt.Errorf("derived tenant base URL must be https (http permitted only for loopback)")
	}
	return strings.TrimRight(base, "/"), nil
}

func writeUnauthorized(w http.ResponseWriter, resourceMetadataURL, scopes, oauthErr, description string) {
	var parts []string
	parts = append(parts, fmt.Sprintf(`resource_metadata=%q`, resourceMetadataURL))
	if scopes != "" {
		parts = append(parts, fmt.Sprintf(`scope=%q`, scopes))
	}
	if oauthErr != "" {
		parts = append(parts, fmt.Sprintf(`error=%q`, oauthErr))
	}
	if description != "" {
		parts = append(parts, fmt.Sprintf(`error_description=%q`, description))
	}
	w.Header().Set("WWW-Authenticate", "Bearer "+strings.Join(parts, ", "))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	_, _ = w.Write([]byte(`{"jsonrpc":"2.0","error":{"code":-32001,"message":"Unauthorized"},"id":null}`))
}

func isLoopbackHost(host string) bool {
	switch strings.ToLower(host) {
	case "localhost", "127.0.0.1", "::1":
		return true
	}
	return false
}

func buildResourceMetadataURL(cfg *config.Config) string {
	u, err := url.Parse(cfg.MCPServerURL)
	if err != nil || u.Host == "" {
		return "/.well-known/oauth-protected-resource/mcp"
	}
	path := strings.TrimRight(u.Path, "/")
	return (&url.URL{
		Scheme: u.Scheme,
		Host:   u.Host,
		Path:   "/.well-known/oauth-protected-resource" + path,
	}).String()
}
