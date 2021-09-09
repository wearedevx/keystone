package jwt

import (
	"fmt"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go/v4"
	"github.com/wearedevx/keystone/api/pkg/models"
	"golang.org/x/xerrors"
)

var salt string

type customClaims struct {
	DeviceUID string `json:"device_uid"`
	jwt.StandardClaims
}

func MakeToken(user models.User, deviceUID string) (string, error) {
	salt := []byte(salt)

	claims := customClaims{
		DeviceUID: deviceUID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: &jwt.Time{
				Time: time.Now().Add(30 * 24 * time.Hour),
			},
			IssuedAt: &jwt.Time{
				Time: time.Now(),
			},
			Issuer:  "keystone",
			Subject: user.UserID,
			// ID:      deviceUID,
		},
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

func VerifyToken(token string) (string, string, error) {
	trimedToken := cleanUpToken(token)

	t, err := jwt.Parse(trimedToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(salt), nil
	})

	if err != nil {
		return "", "", err
	}

	expiredError := &jwt.TokenExpiredError{}

	if t.Valid {
		claims := t.Claims.(jwt.MapClaims)

		userID := claims["sub"].(string)

		if claims["device_uid"] == nil {
			return "", "", fmt.Errorf("Token expired")
		}

		deviceUID := claims["device_uid"].(string)

		return userID, deviceUID, nil
	} else if xerrors.As(err, expiredError) {
		return "", "", fmt.Errorf("Token expired")
	}

	return "", "", fmt.Errorf("Invalid token")
}
