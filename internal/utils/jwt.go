package utils

import (
	"encoding/base64"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	"log"
	"math/rand"
	"os"
	"time"
)

type RefreshToken struct {
	Token       string
	TokenHashed string
}

func loadKey() []byte {
	jwtSecret := os.Getenv("JWT_SECRET_KEY")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET_KEY is empty or not set in .env")
	}

	return []byte(jwtSecret)
}

func GenerateAccessToken(userId string, userIp string) (string, error) {
	jwtToken := jwt.NewWithClaims(
		jwt.SigningMethodHS512,
		jwt.MapClaims{
			"user_id": userId,
			"user_ip": userIp,
			"exp":     time.Now().Add(30 * time.Minute).Unix(),
		})

	signedToken, err := jwtToken.SignedString(loadKey())

	if err != nil {
		return "", err
	}
	return signedToken, nil
}

func GenerateRefreshToken() (RefreshToken, error) {
	byteArr := make([]byte, 32)

	randomSource := rand.NewSource(time.Now().Unix())
	random := rand.New(randomSource)

	if _, err := random.Read(byteArr); err != nil {
		return RefreshToken{}, err
	}

	token := base64.StdEncoding.EncodeToString(byteArr)

	hashedToken, err := bcrypt.GenerateFromPassword([]byte(token), bcrypt.DefaultCost)
	if err != nil {
		return RefreshToken{}, err
	}

	return RefreshToken{token, string(hashedToken)}, nil
}
