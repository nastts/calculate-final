package database

import (
	"errors"

	"github.com/golang-jwt/jwt/v5"
	
)

func ParseToken(tokenString string) (*jwt.Token, error) {
    tokenFromString, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("неверный метод подписи")
		}

		return []byte(HmacSampleSecret), nil
	})

	if err != nil {
		return nil, err
	}
	return tokenFromString, nil
}

func ParseTokenForLogin(tokenString string)(string, error){
	tokenFromString, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("неверный метод подписи")
		}

		return []byte(HmacSampleSecret), nil
	})

	if err != nil {
		return "", err
	}
	
	if claims, ok := tokenFromString.Claims.(jwt.MapClaims); ok {
		result, ok := claims["name"].(string)
		if !ok{
			return "", errors.New("login не найден")
		}
		return result, nil
	} else {
		return "", errors.New("не удалось получить login")

	}
}

