package pkg

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
)

func GenerateSecretJWT(id string) (string, error) {

	var jwtSecret = []byte("your-secret-key")

	claims := jwt.MapClaims{
		"id":  id,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(jwtSecret)

}

func GenerateRefreshJWT(id string) (string, error) {

	var jwtSecret = []byte("your-secret-key")

	claims := jwt.MapClaims{
		"id":  id,
		"exp": time.Now().Add(time.Hour * 24 * 7).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(jwtSecret)

}

func VerifyJWT(tokenString string) (*jwt.Token, error) {
	var jwtSecret = []byte("your-secret-key")
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Check the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return jwtSecret, nil
	})
}
