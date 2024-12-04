package handlers

import (
	"AuthService/internal/models"
	"AuthService/internal/repository"
	"AuthService/internal/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

type TokenRequest struct {
	UserId string `json:"user_id" binding:"required"`
	UserIp string `json:"user_ip" binding:"required"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
	UserIp       string `json:"user_ip" binding:"required"`
	UserId       string `json:"user_id" binding:"required"`
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func GenerateTokens(c *gin.Context) {
	var request TokenRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	accessToken, err := utils.GenerateAccessToken(request.UserId, request.UserIp)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate access token"})
		return
	}

	refreshToken, err := utils.GenerateRefreshToken()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate refresh token"})
		return
	}

	refreshTokenToAdd := models.RefreshToken{
		Id:          uuid.New().String(),
		UserId:      request.UserId,
		TokenHashed: refreshToken.TokenHashed,
		UserIp:      request.UserIp,
		Expires:     time.Now().Add(time.Hour * 48),
	}

	if err = repository.SaveRefreshToken(refreshTokenToAdd); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not save refresh token"})
	}

	response := TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken.Token,
	}

	c.JSON(http.StatusOK, response)

}

func RefreshTokens(c *gin.Context) {
	var request RefreshTokenRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	refreshTokens, err := repository.GetHashedRefreshTokensByUserId(request.UserId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not refresh token"})
		return
	}

	var refreshTokenHashed *models.RefreshToken
	for _, token := range refreshTokens {
		if err = bcrypt.CompareHashAndPassword([]byte(token.TokenHashed), []byte(request.RefreshToken)); err == nil {
			refreshTokenHashed = &token
			break
		}
	}

	if refreshTokenHashed == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
		return
	}

	if time.Now().After(refreshTokenHashed.Expires) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Refresh token expired"})
		return
	}

	if err = bcrypt.CompareHashAndPassword([]byte(refreshTokenHashed.TokenHashed), []byte(request.RefreshToken)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
		return
	}

	if refreshTokenHashed.UserIp != request.UserIp {
		go sendWarningEmail(request.UserId)
	}

	accessToken, err := utils.GenerateAccessToken(request.UserId, request.UserIp)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate access token"})
		return
	}

	refreshToken, err := utils.GenerateRefreshToken()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate refresh token"})
	}

	refreshTokenHashed.TokenHashed = refreshToken.TokenHashed
	refreshTokenHashed.UserIp = request.UserIp
	refreshTokenHashed.Expires = time.Now().Add(time.Hour * 48)

	if err = repository.UpdateRefreshToken(refreshTokenHashed); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not update refresh token"})
		return
	}

	c.JSON(http.StatusOK, TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken.Token,
	})

}
