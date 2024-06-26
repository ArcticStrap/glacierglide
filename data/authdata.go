package data

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/ArcticStrap/glacierglide/jsonresp"
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

func CallJWTAuth(db Datastore, userOnly bool, callback http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u, err := db.GetUser(r.Header.Get("username"))
		if err != nil {
			if userOnly {
				jsonresp.JsonERR(w, http.StatusForbidden, "Invalid token: %s", err)
				return
			}
		} else {
			tokenStr := r.Header.Get("x-jwt-token")
			token, err := ValidateJWT(tokenStr)
			if err != nil || !token.Valid {
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
	return jwt.ParseWithClaims(tokenStr, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {

		// Validate claims
		if claims, ok := token.Claims.(*UserClaims); ok {
			if claims.ExpiresAt.Unix() < time.Now().Unix() {
				return nil, fmt.Errorf("token expired")
			}
		}

		return []byte(os.Getenv("JWTCODE")), nil
	})
}

func ValiateClaims(c jwt.Claims) (jwt.RegisteredClaims, bool) {
	rClaims, ok := c.(jwt.RegisteredClaims)

	return rClaims, ok
}

func GetLoginStatus(tokenStr string, r *http.Request, db Datastore) (string, error) {
	if tokenStr != "" {
		token, err := ValidateJWT(tokenStr)
		if err != nil || !token.Valid {
			return "", err
		}

		claims, ok := token.Claims.(*UserClaims)
		if !ok {
			return "", err
		}

		editor, err := db.GetUsernameFromId(claims.UserID)
		if err != nil {
			return "", err
		}
		return editor, nil
	}
	return strings.Split(r.RemoteAddr, ":")[0], nil
}
