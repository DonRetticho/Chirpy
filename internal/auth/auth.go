package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var ErrNoAuthHeaderIncluded = errors.New("no auth header included in request")

func HashPassword(password string) (string, error) {
	passwordHash, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	return passwordHash, err
}

func CheckPasswordHash(password, hash string) (bool, error) {
	passwordMatch, err := argon2id.ComparePasswordAndHash(password, hash)
	return passwordMatch, err
}

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	claimes := jwt.RegisteredClaims{
		Issuer:    "chirpy-access",
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn).UTC()),
		Subject:   userID.String(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claimes)
	signed, err := token.SignedString([]byte(tokenSecret))
	return signed, err
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(tokenSecret), nil
	})
	if err != nil {
		return uuid.Nil, err
	}
	userIDString, err := token.Claims.GetSubject()
	if err != nil {
		return uuid.Nil, err
	}
	userID, err := uuid.Parse(userIDString)
	if err != nil {
		return uuid.Nil, err
	}
	return userID, err

}

func GetBearerToken(headers http.Header) (string, error) {
	singleHeader := headers.Get("authorization")
	if singleHeader == "" {
		return "", ErrNoAuthHeaderIncluded
	}
	sliceSingleHeader := strings.Split(singleHeader, " ")
	if len(sliceSingleHeader) < 2 || sliceSingleHeader[0] != "Bearer" {
		return "", errors.New("header is malformed")
	}
	token := sliceSingleHeader[1]
	return token, nil
}

func MakeRefreshToken() (string, error) {
	bigByte := make([]byte, 32)
	_, err := rand.Read(bigByte)
	if err != nil {
		return "", err
	}
	encodedString := hex.EncodeToString(bigByte)
	return encodedString, nil
}

func GetAPIKey(headers http.Header) (string, error) {
	singeHeader := headers.Get("authorization")
	if singeHeader == "" {
		return "", ErrNoAuthHeaderIncluded
	}
	sliceSingleHeader := strings.Split(singeHeader, " ")
	if len(sliceSingleHeader) < 2 || sliceSingleHeader[0] != "ApiKey" {
		return "", errors.New("header is malformes")
	}

	apiKey := sliceSingleHeader[1]
	return apiKey, nil
}
