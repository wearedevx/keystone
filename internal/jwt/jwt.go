package jwt

import (
	"fmt"
	"os"
	"strings"

	jwt "github.com/dgrijalva/jwt-go/v4"
	"github.com/wearedevx/keystone/internal/models"
	"github.com/wearedevx/keystone/internal/repo"
	"golang.org/x/xerrors"
)

func MakeToken(user models.User) (string, error) {
	salt := os.Getenv("JWT_SALT")

	claims := jwt.StandardClaims{
		ExpiresAt: jwt.NewTime(24.0 * 60.0 * 60.0),
		IssuedAt:  jwt.NewTime(0.0),
		Issuer:    "keystone",
		Subject:   user.UserID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	fmt.Printf("made token: %v\n", token)

	return token.SignedString(salt)
}

func cleanUpToken(token string) string {
	if strings.HasPrefix(token, "Bearer ") {
		s := strings.Replace(token, "Bearer ", "", 1)

		return strings.Trim(s, " ")
	}

	return strings.Trim(token, " ")
}

func VerifyToken(token string) (string, error) {
	trimedToken := cleanUpToken(token)

	t, err := jwt.Parse(trimedToken, func(token *jwt.Token) (interface{}, error) {
		return os.Getenv("JWT_SALT"), nil
	})

	expiredError := &jwt.TokenExpiredError{}

	if t.Valid {
		Repo := &repo.Repo{}
		Repo.Connect()

		userID := t.Claims.(jwt.StandardClaims).Subject

		return userID, nil
	} else if xerrors.As(err, expiredError) {
		return "", fmt.Errorf("Token expired")
	}

	return "", fmt.Errorf("Invalid token")
}
