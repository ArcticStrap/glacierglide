package data

import (
	"fmt"
	"net/http"
	"os"

	"github.com/ChaosIsFramecode/horinezumi/jsonresp"
	"github.com/golang-jwt/jwt/v5"
)

func CreateJWT(u *User) (string, error) {
	claims := &jwt.MapClaims{
		"expiresAt": 15000,
		"userId":    u.UserId,
	}

	jwtcode := os.Getenv("JWTCODE")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(jwtcode))
}

func CallJWTAuth(db Datastore, callback http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenStr := r.Header.Get("x-jwt-token")
		token, err := ValidateJWT(tokenStr)
		if err != nil || !token.Valid {
			jsonresp.JsonERR(w, http.StatusForbidden, "Invalid token: %s", err)
			return
		}

		u, err := db.GetUser(r.Header.Get("username"))
		if err != nil {
			jsonresp.JsonERR(w, http.StatusForbidden, "Invalid token: %s", err)
		}

		claims := token.Claims.(jwt.MapClaims)
		if int64(claims["userId"].(float64)) != u.UserId {
			jsonresp.JsonERR(w, http.StatusForbidden, "Invalid token: %s", err)
		}

		callback(w, r)
	}
}

func ValidateJWT(tokenStr string) (*jwt.Token, error) {
	jwtcode := os.Getenv("JWTCODE")

	return jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		// Valide algorithm
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(jwtcode), nil
	})
}
