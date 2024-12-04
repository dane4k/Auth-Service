package repository

import (
	"AuthService/db"
	"AuthService/internal/models"
	"github.com/google/uuid"
)

func SaveRefreshToken(token models.RefreshToken) error {
	return db.DB.Create(&token).Error
}

func GetHashedRefreshTokensByUserId(userId string) ([]models.RefreshToken, error) {
	var refreshTokens []models.RefreshToken
	if err := db.DB.Where("user_id = ?", userId).Find(&refreshTokens).Error; err != nil {
		return nil, err
	}
	return refreshTokens, nil
}

func UpdateRefreshToken(refreshToken *models.RefreshToken) error {
	return db.DB.Model(&models.RefreshToken{}).
		Where("id = ?", refreshToken.Id).
		Updates(map[string]interface{}{
			"token_hashed": refreshToken.TokenHashed,
			"user_ip":      refreshToken.UserIp,
			"expires":      refreshToken.Expires}).Error
}

// AddUser useless, for tests
func AddUser(email string) error {
	user := models.User{
		Id:    uuid.New().String(),
		Email: email,
	}

	if err := db.DB.Create(&user).Error; err != nil {
		return err
	}

	return nil
}

func GetUserIdByEmail(email string) (string, error) {
	var user models.User
	if err := db.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return "", err
	}
	return user.Id, nil
}
