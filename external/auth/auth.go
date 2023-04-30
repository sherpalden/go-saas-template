package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

const (
	ACCESS_TOKEN_VALIDITY_PERIOD  time.Duration = 60 * time.Minute
	REFRESH_TOKEN_VALIDITY_PERIOD time.Duration = 7 * 24 * 60 * time.Minute
)

type LoginCredential struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type AuthTokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type AuthClaims struct {
	*jwt.StandardClaims
	Custom map[string]interface{} `json:"custom"`
}

type AuthClient struct {
	AccessTokenSecret  string
	RefreshTokenSecret string
}

func NewAuthClient(
	accessTokenSecret string,
	refreshTokenSecret string,
) AuthClient {
	return AuthClient{
		AccessTokenSecret:  accessTokenSecret,
		RefreshTokenSecret: refreshTokenSecret,
	}
}

func createToken(signKey string, customClaims map[string]interface{}, validityPeriod time.Duration) (string, error) {
	claims := &AuthClaims{
		&jwt.StandardClaims{
			ExpiresAt: time.Now().Add(validityPeriod).Unix(),
			IssuedAt:  time.Now().Unix(),
			NotBefore: time.Now().Unix(),
			Issuer:    "erp",
			Subject:   "auth",
		},
		customClaims,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(signKey))
}

func (ac *AuthClient) GenerateTokens(customClaims map[string]interface{}) (AuthTokens, error) {
	accessToken, err := createToken(ac.AccessTokenSecret, customClaims, ACCESS_TOKEN_VALIDITY_PERIOD)
	if err != nil {
		return AuthTokens{}, err
	}

	refreshToken, err := createToken(ac.RefreshTokenSecret, customClaims, REFRESH_TOKEN_VALIDITY_PERIOD)
	if err != nil {
		return AuthTokens{}, err
	}

	return AuthTokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (ac *AuthClient) VerifyAccessToken(tokenString string) (claims map[string]interface{}, err error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(ac.AccessTokenSecret), nil
	})
	if err != nil {
		return claims, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if _, ok := claims["custom"]; ok {
			return claims["custom"].(map[string]interface{}), nil
		}
		return map[string]interface{}{}, nil
	} else {
		return map[string]interface{}{}, nil
	}
}

func (ac *AuthClient) VerifyRefreshToken(tokenString string) (claims map[string]interface{}, err error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(ac.RefreshTokenSecret), nil
	})
	if err != nil {
		return claims, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if _, ok := claims["custom"]; ok {
			return claims["custom"].(map[string]interface{}), nil
		}
		return map[string]interface{}{}, nil
	} else {
		return map[string]interface{}{}, nil
	}
}
