package services

import (
	"Alarm/internal/web/models"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Auth struct {
	db    *models.Database
	cache *models.Cache
	cfg   map[string]interface{}
}

func NewAuth(cfg map[string]interface{}) *Auth {
	return &Auth{
		db:    cfg["db"].(*models.Database),
		cache: cfg["cache"].(*models.Cache),
		cfg:   cfg,
	}
}

// RefreshToken 用于刷新token，可以创建新的token或者延长现有token的过期时间。
//
// 如果 tokenStr 为空字符串，则使用私钥创建新的token。
func (svc *Auth) RefreshToken(user *models.User, tokenStr string, validSeconds int) (string, error) {
	var err error
	if tokenStr == "" {
		tokenStr, err = generateToken(user, validSeconds, svc.cfg["privateKey"])
		if err != nil {
			return "", err
		}
	}
	return tokenStr, nil
}

func (svc *Auth) ValidateToken(tokenStr string) (*jwt.MapClaims, error) {
	// 读取公钥
	publicKey := svc.cfg["publicKey"]

	// 解析令牌
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		// 验证签名方法是否为RS256
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("invalid signing method: %v", token.Header["alg"])
		}

		return publicKey, nil
	})
	if err != nil {
		return nil, err
	}

	// 验证令牌是否有效
	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	// 解析声明
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid claims")
	}

	// 验证过期时间
	expirationTime := time.Unix(int64(claims["exp"].(float64)), 0)
	if time.Now().UTC().After(expirationTime) {
		return nil, fmt.Errorf("token has expired")
	}

	return &claims, nil
}

func (svc *Auth) DeleteToken(token string) error {
	redisKey := fmt.Sprintf("token:%s", token)
	err := svc.cache.Client.Del(redisKey).Err()
	if err != nil {
		return err
	}
	return nil
}

func (svc *Auth) VerifyPassword(username string, password string) (int, error) {
	user := &models.User{
		Username: username,
	}
	has, err := svc.db.Engine.Cols("id", "password").Get(user)
	if err != nil {
		return 0, err
	}
	if !has {
		return 0, errors.New("user not found")
	}
	return user.ID, nil
}

func generateToken(user *models.User, validSeconds int, privateKey interface{}) (string, error) {
	claims := &jwt.MapClaims{
		"userID": user.ID,
		"exp":    time.Now().Add(time.Second * time.Duration(validSeconds)).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	tokenStr, err := token.SignedString(privateKey)
	if err != nil {
		return "", err
	}

	return tokenStr, nil
}
