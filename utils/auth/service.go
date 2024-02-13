package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	jwt "github.com/golang-jwt/jwt/v4"
)

var (
	jwtKey     []byte
	jwtKeyName JWTKeyName = "tkn"
)

type JWTKeyName string

type claimsData struct {
	UserID          int64  `json:"uid"`
	CompositeUserId string `json:"cid"`
	Platform        string `json:"ptf"`
	Scope           string `json:"scope"`
	Role            int    `json:"rle"`
	LID             int    `json:"lid"`

	jwt.StandardClaims
}

func Init(key string) {
	jwtKey = []byte(key)
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "application/json")

		token, err := TokenFromRequest(string(jwtKey), r)
		if err != nil {
			if token == nil {
				w.WriteHeader(http.StatusForbidden)
				return
			}
			if token.Valid {
				log.Println("You look nice today")
			} else if errors.Is(err, jwt.ErrTokenMalformed) {
				log.Println("That's not even a token")
				w.WriteHeader(http.StatusForbidden)
				return
			} else if errors.Is(err, jwt.ErrTokenExpired) || errors.Is(err, jwt.ErrTokenNotValidYet) {
				// Token is either expired or not active yet
				log.Println("Timing is everything")
				w.WriteHeader(http.StatusForbidden)
				return
			} else {
				log.Println("Couldn't handle this token:", err)
				w.WriteHeader(http.StatusForbidden)
				return
			}
		}

		ctx := context.WithValue(r.Context(), jwtKeyName, token)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Get the token from the http header
func TokenFromRequest(signKey string, r *http.Request) (*jwt.Token, error) {
	const bearerPfx = "Bearer "
	const TokenHeader = "Authorization"

	authHeader := r.Header.Get(TokenHeader)
	if len(authHeader) == 0 {
		return nil, errors.New("missing header")
	}

	if !strings.HasPrefix(authHeader, bearerPfx) {
		return nil, errors.New("invalid token prefix")
	}

	tkn := strings.TrimPrefix(authHeader, bearerPfx)
	return jwt.Parse(tkn, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(signKey), nil
	})
}

func ClaimsFromCtx(ctx context.Context) (*claimsData, error) {
	if ctx.Value(jwtKeyName) == nil {
		return nil, errors.New("token not found")
	}

	token := ctx.Value(jwtKeyName).(*jwt.Token)
	if token == nil {
		return nil, errors.New("token not found")
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	cd, err := json.Marshal(token.Claims)
	if err != nil {
		return nil, errors.New("invalid token")
	}

	c := new(claimsData)
	if err := json.Unmarshal(cd, c); err != nil {
		return nil, err
	}

	return c, nil
}
