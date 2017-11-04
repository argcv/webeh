package webeh

import (
	"net/http"
	"strconv"
	"strings"
)

type Adapter func(http.Handler) http.Handler

var AccessControlAllowMethods = []string{
	"POST",
	"GET",
	"OPTIONS",
	"PUT",
	"PATCH",
	"DELETE",
}

var AccessControlAllowHeaders = []string{
	"Content-Type",
	"Content-Length",
	"Accept-Encoding",
	"Content-Range",
	"Content-Disposition",
	"Authorization",
	"X-Requested-With",
}

var AccessControlAllowCredentials = true

func Adapt(h http.Handler, adapters ...Adapter) http.Handler {
	for k := len(adapters) - 1; k >= 0; k-- {
		h = adapters[k](h)
	}
	return h
}

func CORSMiddleware() Adapter {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Methods", strings.Join(AccessControlAllowMethods, ","))
			w.Header().Set("Access-Control-Allow-Headers", strings.Join(AccessControlAllowHeaders, ","))
			// Since we need to support cross-domain cookies, we must support XHR requests
			// with credentials, so the Access-Control-Allow-Credentials header is required
			// See https://developer.mozilla.org/en-US/docs/Web/HTTP/Access_control_CORS.
			w.Header().Set("Access-Control-Allow-Credentials", strconv.FormatBool(AccessControlAllowCredentials))
			w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))

			if r.Method == "OPTIONS" {
				return
			}

			h.ServeHTTP(w, r)
		})
	}
}
