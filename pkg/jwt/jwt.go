package jwt

import (
	"golang.org/x/oauth2/jwt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var signingKey = []byte("secret")

type UserInfo struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

type UserClaims struct {
	User *UserInfo `json:"user"`
	jwt.StandardClaims
}

func GetToken(id, username string) (string, error) {
	now := time.Now()
	claims := &UserClaims{
		User: &UserInfo{
			ID:       id,
			Username: username,
		},
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  now.Unix(),
			ExpiresAt: now.AddDate(0, 1, 0).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(signingKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func ParseToken(tokenString string) (*UserInfo, error) {
	token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		return signingKey, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*UserClaims); ok && token.Valid {
		return claims.User, nil
	}
	return nil, err
}
