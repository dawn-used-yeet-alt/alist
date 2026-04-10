package db

import (
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/alist-org/alist/v3/internal/model"
)

func GenerateToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func CreateShareToken(path, label string, expiresAt *time.Time) (*model.ShareToken, error) {
	token, err := GenerateToken()
	if err != nil {
		return nil, err
	}
	st := &model.ShareToken{
		Token:     token,
		Path:      path,
		Label:     label,
		ExpiresAt: expiresAt,
	}
	return st, db.Create(st).Error
}

func GetShareToken(token string) (*model.ShareToken, error) {
	var st model.ShareToken
	err := db.Where("token = ?", token).First(&st).Error
	return &st, err
}

func ListShareTokens() ([]model.ShareToken, error) {
	var tokens []model.ShareToken
	err := db.Find(&tokens).Error
	return tokens, err
}

func DeleteShareToken(id uint) error {
	return db.Delete(&model.ShareToken{}, id).Error
}
