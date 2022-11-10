package jwt

import (
	"errors"
	"strings"
	"time"

	jwt "github.com/golang-jwt/jwt/v4"

	"github.com/wearedevx/keystone/api/internal/utils"
	"github.com/wearedevx/keystone/api/pkg/models"
)

var salt string

var (
	ErrorInvalidToken = errors.New("invalid token")
	ErrorTokenExpired = errors.New("token expired")
)

type customClaims struct {
	DeviceUID string `json:"device_uid"`
	jwt.StandardClaims
}

func createRefreshToken() (string, error) {
	return utils.RandomString(64)
}

func MakeToken(user models.User, deviceUID string, when time.Time) (string, string, error) {
	salt := []byte(salt)

	claims := customClaims{
		DeviceUID: deviceUID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: when.Add(24 * time.Hour).Unix(),
			IssuedAt:  when.Unix(),
			Issuer:    "keystone",
			Subject:   user.UserID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	refreshToken, err := createRefreshToken()
	if err != nil {
		return "", "", err
	}

	signed, err := token.SignedString(salt)
	if err != nil {
		return "", "", err
	}

	return signed, refreshToken, nil
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

	t, err := jwt.Parse(
		trimedToken,
		func(token *jwt.Token) (interface{}, error) {
			return []byte(salt), nil
		},
	)
	if err != nil {
		if validationError, ok := err.(jwt.ValidationError); ok {
			if validationError.Errors&jwt.ValidationErrorExpired == jwt.ValidationErrorExpired {
				return "", "", ErrorTokenExpired
			}
		} else {
			return "", "", err
		}
	}

	if t.Valid {
		claims := t.Claims.(jwt.MapClaims)

		userID := claims["sub"].(string)

		if claims["device_uid"] == nil {
			return "", "", ErrorTokenExpired
		}

		deviceUID := claims["device_uid"].(string)

		return userID, deviceUID, nil
	}

	return "", "", ErrorInvalidToken
}
