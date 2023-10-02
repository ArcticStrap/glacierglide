package data

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/ChaosIsFramecode/horinezumi/jsonresp"
	"github.com/golang-jwt/jwt/v5"
)

type UserClaims struct {
	UserID int64
	jwt.RegisteredClaims
}

func CreateJWT(u *User) (string, error) {
	claims := &UserClaims{
		UserID: u.UserId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(time.Hour * time.Duration(24))),
		},
	}

	jwtcode := os.Getenv("JWTCODE")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(jwtcode))
}

func CallJWTAuth(db Datastore, callback http.HandlerFunc, userOnly bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if userOnly {
			tokenStr := r.Header.Get("x-jwt-token")
			token, err := ValidateJWT(tokenStr)
			if err != nil || !token.Valid {
				jsonresp.JsonERR(w, http.StatusForbidden, "Invalid token: %s", err)
				return
			}

			u, err := db.GetUser(r.Header.Get("username"))
			if err != nil {
				jsonresp.JsonERR(w, http.StatusForbidden, "Invalid token: %s", err)
				return
			}

			claims := token.Claims.(*UserClaims)
			if claims.UserID != u.UserId {
				jsonresp.JsonERR(w, http.StatusForbidden, "Invalid token: %s", err)
				return
			}
		}

		callback(w, r)
	}
}

func ValidateJWT(tokenStr string) (*jwt.Token, error) {
	return jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		// Validate algorithm
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// Validate claims
		if claims, ok := token.Claims.(UserClaims); ok {
			if claims.ExpiresAt.Unix() < time.Now().Unix() {
				return nil, fmt.Errorf("token expired")
			}
		}

		return []byte(os.Getenv("JWTCODE")), nil
	})
}
