package middleware

import (
	"net"
	"net/http"
)

func IPChecker(trustedSubnet string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			if trustedSubnet == "" {
				http.Error(res, http.StatusText(http.StatusForbidden), http.StatusForbidden)
				return
			}

			_, ipNet, err := net.ParseCIDR(trustedSubnet)
			if err != nil {
				http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}

			xRealIP := req.Header.Get("X-Real-IP")
			if xRealIP == "" {
				http.Error(res, http.StatusText(http.StatusForbidden), http.StatusForbidden)
				return
			}

			realIP := net.ParseIP(xRealIP)
			if !ipNet.Contains(realIP) {
				http.Error(res, http.StatusText(http.StatusForbidden), http.StatusForbidden)
				return
			}

			next.ServeHTTP(res, req)
		})
	}
}
