package repositories

import (
	"context"
	"log"
	"net/http"
	"sample-backend-go/internal/domain/entities/config"

	"github.com/Nerzal/gocloak/v13"
)

type LoginError struct {
	Code    int
	Message string
}

func (e *LoginError) Error() string {
	return e.Message
}

func LoginToKeycloakUsername(username, password string) (string, string, error) {
	cfg, _ := config.NewConfig()

	client := gocloak.NewClient(cfg.KeycloakApiUrl)

	realm := cfg.KeycloakRealm
	clientID := cfg.KeycloakClientId
	clientSecret := cfg.KeycloakClientSecret

	token, err := client.Login(context.Background(), clientID, clientSecret, realm, username, password)
	if err != nil {
		apiErr, ok := err.(*gocloak.APIError)
		if ok {
			if apiErr.Code == http.StatusUnauthorized {
				log.Printf("認証失敗: ユーザーネームかパスワードが間違っています: %v", err)
				return "", "", &LoginError{Code: apiErr.Code, Message: "Username or password is incorrect"}
			}
			log.Printf("Keycloakとの通信に失敗: %v", apiErr.Message)
			return "", "", &LoginError{Code: apiErr.Code, Message: apiErr.Message}
		} else {
			log.Printf("予期せぬエラーが発生: %v", err)
			return "", "", &LoginError{Code: http.StatusInternalServerError, Message: "Internal server error"}
		}
	}
	return token.AccessToken, token.RefreshToken, nil
}
func LoginToKeycloakEmail(email, password string) (string, string, error) {
	cfg, _ := config.NewConfig()

	client := gocloak.NewClient(cfg.KeycloakApiUrl)

	realm := cfg.KeycloakRealm
	clientID := cfg.KeycloakClientId
	clientSecret := cfg.KeycloakClientSecret

	token, err := client.Login(context.Background(), clientID, clientSecret, realm, email, password)
	if err != nil {
		apiErr, ok := err.(*gocloak.APIError)
		if ok {
			if apiErr.Code == http.StatusUnauthorized {
				log.Printf("認証失敗: メールアドレスかパスワードが間違っています: %v", err)
				return "", "", &LoginError{Code: apiErr.Code, Message: "Email or password is incorrect"}
			}
			log.Printf("Keycloakとの通信に失敗: %v", apiErr.Message)
			return "", "", &LoginError{Code: apiErr.Code, Message: apiErr.Message}
		} else {
			log.Printf("予期せぬエラーが発生: %v", err)
			return "", "", &LoginError{Code: http.StatusInternalServerError, Message: "Internal server error"}
		}
	}
	return token.AccessToken, token.RefreshToken, nil
}

func RefreshToken(refreshToken string) (string, string, error) {
	cfg, _ := config.NewConfig()

	client := gocloak.NewClient(cfg.KeycloakApiUrl)

	realm := cfg.KeycloakRealm
	clientID := cfg.KeycloakClientId
	clientSecret := cfg.KeycloakClientSecret

	token, err := client.RefreshToken(context.Background(), refreshToken, clientID, clientSecret, realm)
	if err != nil {
		apiErr, ok := err.(*gocloak.APIError)
		if ok {
			log.Printf("Keycloakとの通信に失敗: %v", apiErr.Message)
			return "", "", &LoginError{Code: apiErr.Code, Message: apiErr.Message}
		} else {
			log.Printf("予期せぬエラーが発生: %v", err)
			return "", "", &LoginError{Code: http.StatusInternalServerError, Message: "Internal server error"}
		}
	}
	return token.AccessToken, token.RefreshToken, nil
}
