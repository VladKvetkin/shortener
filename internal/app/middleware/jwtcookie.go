package middleware

import (
	"context"
	"net/http"

	"github.com/VladKvetkin/shortener/internal/app/auth"
)

const TokenCookieName = "token"

type UserIDKey struct{}

// JWTCookie - функция, которая проверяет наличие JWT-токена, и если его нет, то генерирует его и записывает в куки.
func JWTCookie(next http.Handler) http.Handler {
	return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		var token string

		tokenCookie, err := req.Cookie(TokenCookieName)

		if err != nil {
			token, err = auth.BuildJWTToken()
			if err != nil {
				http.Error(resp, "Cannot build JWT for user", http.StatusBadRequest)
				return
			}
		} else {
			token = tokenCookie.Value
			_, err := auth.GetUserID(token)
			if err != nil {
				token, err = auth.BuildJWTToken()
				if err != nil {
					http.Error(resp, "Cannot build JWT for user", http.StatusBadRequest)
					return
				}
			}
		}

		userID, err := auth.GetUserID(token)
		if err != nil {
			http.Error(resp, "Invalid JWT", http.StatusBadRequest)
			return
		}

		req = req.WithContext(context.WithValue(req.Context(), UserIDKey{}, userID))

		http.SetCookie(resp, &http.Cookie{
			Name:  TokenCookieName,
			Value: token,
			Path:  "/",
		})

		next.ServeHTTP(resp, req)
	})
}
