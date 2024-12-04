package tests

import (
	"AuthService/db"
	"AuthService/internal/handlers"
	"AuthService/internal/repository"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var refreshToken, currentUserId string

func TestGenerateTokens(t *testing.T) {

	if err := godotenv.Load("../../.env"); err != nil {
		fmt.Println(err)
	}

	db.InitDB()

	router := gin.Default()

	email := "fake@gmail.com"

	err := repository.AddUser(email)
	if err != nil {
		log.Println(err)
	}

	userId, err := repository.GetUserIdByEmail(email)
	if err != nil {
		t.Log(err)
	}

	currentUserId = userId

	requestData := handlers.TokenRequest{
		UserId: userId,
		UserIp: "127.0.0.1",
	}

	requestBodyValid, err := json.Marshal(requestData)
	if err != nil {
		t.Log(err)
	}

	request, _ := http.NewRequest("POST", "/generate_tokens", bytes.NewReader(requestBodyValid))

	recorder := httptest.NewRecorder()

	router.POST("/generate_tokens", handlers.GenerateTokens)

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Errorf("got status code: %d, want: %d", recorder.Code, http.StatusOK)
	} else if !strings.Contains(recorder.Body.String(), "access_token") || !strings.Contains(recorder.Body.String(), "refresh_token") {
		t.Errorf("expected response body to contain access_token, got: %s", recorder.Body.String())
	}

	var response handlers.TokenResponse
	if err = json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatal(err)
	}

	if response.AccessToken == "" || response.RefreshToken == "" {
		t.Errorf("want access_token and refresh_token, but got: %v", response)
	}

	refreshToken = response.RefreshToken
}

func TestRefreshTokens(t *testing.T) {
	if err := godotenv.Load("../../.env"); err != nil {
		fmt.Println(err)
	}

	db.InitDB()

	router := gin.Default()

	requestDataValid := handlers.RefreshTokenRequest{
		RefreshToken: refreshToken,
		UserIp:       "127.0.0.1",
		UserId:       currentUserId,
	}

	requestBodyValid, err := json.Marshal(requestDataValid)
	if err != nil {
		t.Log(err)
	}

	ValidRequest, _ := http.NewRequest("POST", "/refresh_tokens", bytes.NewReader(requestBodyValid))

	recorder := httptest.NewRecorder()

	router.POST("/refresh_tokens", handlers.RefreshTokens)

	router.ServeHTTP(recorder, ValidRequest)

	if recorder.Code != http.StatusOK {
		t.Errorf("got status code: %d, want: %d", recorder.Code, http.StatusOK)
	}
	var response handlers.TokenResponse
	if err = json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatal(err)
	}
	if response.AccessToken == "" || response.RefreshToken == "" {
		t.Errorf("want access_token and refresh_token, but got: %v", response)
	}
}
