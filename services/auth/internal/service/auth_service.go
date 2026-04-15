package service

import (
	"errors"
	"time"
	"bytes"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"simkeu/service-auth/internal/repository"

	"encoding/json"
	"net/http"
	"fmt"
	"log"
	"io"
)

type AuthService struct {
	Repo      *repository.UserRepository
	JWTSecret string
	DebiturURL string
}

func (s *AuthService) Register(email, password string) error {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	err = s.Repo.Create(email, string(hashed))
	if err != nil {
		return err
	}

	id, _, err := s.Repo.FindByEmail(email)
	if err != nil {
		return err
	}

	go s.CreateDebitur(id, email)

	return nil
}

func (s *AuthService) Login(email, password string) (string, error) {
	id, storedPassword, err := s.Repo.FindByEmail(email)
	if err != nil {
		return "", errors.New("invalid email or password")
	}

	err = bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(password))
	if err != nil {
		return "", errors.New("invalid email or password")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": id,
		"email":   email,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	})

	return token.SignedString([]byte(s.JWTSecret))
}

func (s *AuthService) ValidateToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.JWTSecret), nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid claims")
	}

	return claims, nil
}

func (s *AuthService) GetDebiturProfile(userID int, token string) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/api/debitur/%d", s.DebiturURL, userID)

	client := http.Client{
		Timeout: 3 * time.Second,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// 🔥 DEBUG
	log.Println("STATUS CODE:", resp.StatusCode)

	bodyBytes, _ := io.ReadAll(resp.Body)
	log.Println("RESPONSE BODY:", string(bodyBytes))

	// ❗ penting: decode ulang dari bytes
	var result map[string]interface{}
	err = json.Unmarshal(bodyBytes, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *AuthService) CreateDebitur(id int, name string) {
	url := fmt.Sprintf("%s/api/debitur", s.DebiturURL)

	payload := map[string]interface{}{
		"id":   id,
		"name": name,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		log.Println("failed to marshal payload:", err)
		return
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Println("failed to create request:", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")

	client := http.Client{Timeout: 3 * time.Second}

	resp, err := client.Do(req)
	if err != nil {
		log.Println("failed to call debitur service:", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	log.Println("CreateDebitur response:", string(body))
}
