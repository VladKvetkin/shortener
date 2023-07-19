package middleware

import (
	"context"
	"net/http"

	"github.com/VladKvetkin/shortener/internal/app/auth"
)

const TokenCookieName = "token"

type UserIDKey struct{}

func JWTCookie(next http.Handler) http.Handler {
	return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		var token string

		jwtTokenBuilder := auth.JwtTokenBuilder{}

		tokenCookie, err := req.Cookie(TokenCookieName)

		if err != nil {
			token, err = jwtTokenBuilder.BuildJWTToken()
			if err != nil {
				http.Error(resp, "Cannot build JWT for user", http.StatusBadRequest)
				return
			}
		} else {
			token = tokenCookie.Value
			userID, err := jwtTokenBuilder.GetUserID(token)
			if err != nil {
				token, err = jwtTokenBuilder.BuildJWTToken()
				if err != nil {
					http.Error(resp, "Cannot build JWT for user", http.StatusBadRequest)
					return
				}
			} else {
				req = req.WithContext(context.WithValue(req.Context(), UserIDKey{}, userID))
			}
		}

		http.SetCookie(resp, &http.Cookie{
			Name:  TokenCookieName,
			Value: token,
			Path:  "/",
		})

		next.ServeHTTP(resp, req)
	})
}
