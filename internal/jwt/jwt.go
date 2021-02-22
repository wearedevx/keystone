package jwt

import (
	"fmt"
	"os"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go/v4"
	"github.com/wearedevx/keystone/internal/models"
	"github.com/wearedevx/keystone/internal/repo"
	"golang.org/x/xerrors"
)

func MakeToken(user models.User) (string, error) {
	salt := os.Getenv("JWT_SALT")

	claims := jwt.StandardClaims{
		ExpiresAt: jwt.Time(time.Now().Add(time.Hour * 24)),
		IssuedAt:  jwt.Time(time.Now()),
		Issuer:    "keystone",
		Subject:   user.UserID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

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
