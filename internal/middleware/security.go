package middleware

import "github.com/gogf/gf/v2/net/ghttp"

// SecurityHeaders sets common security headers on all responses.
func SecurityHeaders(r *ghttp.Request) {
	r.Response.Header().Set("X-Content-Type-Options", "nosniff")
	r.Response.Header().Set("X-Frame-Options", "DENY")
	r.Response.Header().Set("X-XSS-Protection", "1; mode=block")
	r.Response.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
	r.Response.Header().Set("Content-Security-Policy",
		"default-src 'self'; script-src 'self' 'unsafe-inline' 'unsafe-eval'; style-src 'self' 'unsafe-inline'; img-src 'self' data:;")
	r.Middleware.Next()
}
