package api

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

type UserClaims struct {
	ID string `json:"id"`
	jwt.StandardClaims
}

type ctxKey int

const (
	userIDCtxKey ctxKey = iota
)

func (s *Server) AuthenticateMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		tokenString := r.Header.Get("Authorization")
		splitToken := strings.Split(tokenString, "Bearer ")
		if len(splitToken) < 2 {
			log.Printf("tocken is not found in header: %s", tokenString)
			rw.WriteHeader(401)
			return
		}
		tokenString = splitToken[1]

		token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(token *jwt.Token) (interface{}, error) { return s.key, nil })
		if err != nil || !token.Valid {
			log.Printf("token is invalid %s: %v", tokenString, err)
			rw.WriteHeader(401)
			return
		}

		claims, ok := token.Claims.(*UserClaims)
		if !ok {
			log.Printf("token is invalid: %s", tokenString)
			rw.WriteHeader(401)
			return
		}

		ctx = context.WithValue(ctx, userIDCtxKey, claims.ID)

		r = r.WithContext(ctx)

		next.ServeHTTP(rw, r)
	})
}
