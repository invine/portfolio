package ports

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

type UserClaims struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Login string `json:"login"`
	Name  string `json:"name"`
	jwt.StandardClaims
}

type User struct {
	ID uuid.UUID
}

type ctxKey int

const (
	userCtxKey ctxKey = iota
)

func (s *Server) AuthenticateMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		tokenString := r.Header.Get("Authorization")
		splitToken := strings.Split(tokenString, "Bearer ")
		if len(splitToken) < 2 {
			log.Printf("token is not found in header: %s", tokenString)
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

		id, err := uuid.Parse(claims.ID)
		if err != nil {
			log.Printf("user id %s is invalid: %v", claims.ID, err)
			rw.WriteHeader(401)
			return
		}
		ctx = context.WithValue(ctx, userCtxKey, User{ID: id})

		r = r.WithContext(ctx)

		next.ServeHTTP(rw, r)
	})
}

func (s *Server) UserSignInHandler(rw http.ResponseWriter, r *http.Request) {
	bytes, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("failed sing in: %v", err)
		rw.WriteHeader(400)
		return
	}

	type credentials struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}

	creds := new(credentials)
	if err := json.Unmarshal(bytes, creds); err != nil {
		log.Printf("failed sign in: %v", err)
		rw.WriteHeader(400)
		return
	}

	ctx := context.Background()

	u, err := s.userSvc.AuthenticateUser(ctx, creds.Login, creds.Password)
	if err != nil {
		log.Printf("failed sign in: %v", err)
		rw.WriteHeader(401)
		return
	}

	// Create the Claims
	claims := UserClaims{
		u.ID().String(),
		u.Email(),
		u.Login(),
		u.Name(),
		jwt.StandardClaims{
			ExpiresAt: time.Now().Unix() + 15000,
			Issuer:    "portfolio",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(s.key)
	if err != nil {
		log.Printf("failed sign in: %v", err)
		rw.WriteHeader(500)
		return
	}

	if _, err := rw.Write([]byte(ss)); err != nil {
		log.Printf("failed sign in: %v", err)
	}
}

func (s *Server) UserSignUpHandler(rw http.ResponseWriter, r *http.Request) {
	bytes, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("failed sign up: %v", err)
		rw.WriteHeader(400)
		return
	}

	type userModel struct {
		Email    string `json:"email"`
		Login    string `json:"login"`
		Password string `json:"password"`
		Name     string `json:"name"`
	}

	um := new(userModel)
	if err := json.Unmarshal(bytes, um); err != nil {
		log.Printf("failed sign up: %v", err)
		rw.WriteHeader(400)
		return
	}

	ctx := context.Background()
	if err := s.userSvc.CreateUser(ctx, um.Email, um.Login, um.Password, um.Name); err != nil {
		log.Printf("failed sign up: %v", err)
		rw.WriteHeader(500)
		return
	}

	rw.WriteHeader(200)
}

func UserFromCtx(ctx context.Context) (User, error) {
	u, ok := ctx.Value(userCtxKey).(User)
	if !ok {
		return User{}, fmt.Errorf("no user in context")
	}

	return u, nil
}
