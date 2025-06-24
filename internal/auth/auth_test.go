package auth

import (
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestCheckPasswordHash(t *testing.T) {
	password1 := "correctPassword123!"
	password2 := "anotherPassword456!"
	hash1, _ := HashPassword(password1)
	hash2, _ := HashPassword(password2)

	tests := []struct {
		name     string
		password string
		hash     string
		wantErr  bool
	}{
		{
			name:     "Correct password",
			password: password1,
			hash:     hash1,
			wantErr:  false,
		},
		{
			name:     "Incorrect password",
			password: "wrongPassword",
			hash:     hash1,
			wantErr:  true,
		},
		{
			name:     "Password doesn't match different hash",
			password: password1,
			hash:     hash2,
			wantErr:  true,
		},
		{
			name:     "Empty password",
			password: "",
			hash:     hash1,
			wantErr:  true,
		},
		{
			name:     "Invalid hash",
			password: password1,
			hash:     "invalidhash",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := CheckPasswordHash(tt.password, tt.hash)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckPasswordHash() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateJWT(t *testing.T) {
	uuid_valid := uuid.New()
	tokenSecret_valid := "superSecretKey123!"
	jwt_valid, _ := MakeJWT(uuid_valid, tokenSecret_valid, 24*time.Hour)
	uuid_expired := uuid.New()
	tokenSecret_expired := "anotherSecretKey456!"
	jwt_expired, _ := MakeJWT(uuid_expired, tokenSecret_expired, 1*time.Nanosecond)

	tests := []struct {
		name           string
		tokenString    string
		tokenSecret    string
		expectedUserId uuid.UUID
		wantErr        bool
	}{
		{
			name:           "Valid JWT",
			tokenString:    jwt_valid,
			tokenSecret:    tokenSecret_valid,
			expectedUserId: uuid_valid,
			wantErr:        false,
		},
		{
			name:           "Expired JWT",
			tokenString:    jwt_expired,
			tokenSecret:    tokenSecret_expired,
			expectedUserId: uuid.Nil,
			wantErr:        true,
		},
		{
			name:           "JWT signed with wrong secret",
			tokenString:    jwt_valid,
			tokenSecret:    "wrongSecretKey789!",
			expectedUserId: uuid.Nil,
			wantErr:        true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotUUID, err := ValidateJWT(tt.tokenString, tt.tokenSecret)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateJWT() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotUUID != tt.expectedUserId {
				t.Errorf("ValidateJWT() gotUUID = %v, want %v", gotUUID, tt.expectedUserId)
			}
		})
	}
}

func TestGetBearerToken(t *testing.T) {
	tests := []struct {
		name          string
		headers       map[string]string
		expectedToken string
		wantErr       bool
	}{
		{
			name:          "Valid Bearer token",
			headers:       map[string]string{"Authorization": "Bearer validToken123"},
			expectedToken: "validToken123",
			wantErr:       false,
		},
		{
			name:          "Missing Authorization header",
			headers:       map[string]string{},
			expectedToken: "",
			wantErr:       true,
		},
		{
			name:          "Invalid Authorization format",
			headers:       map[string]string{"Authorization": "BearervalidToken123"},
			expectedToken: "",
			wantErr:       true,
		},
		{
			name:          "Empty Bearer token",
			headers:       map[string]string{"Authorization": "Bearer "},
			expectedToken: "",
			wantErr:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			headers := make(http.Header)
			for k, v := range tt.headers {
				headers.Set(k, v)
			}
			gotToken, err := GetBearerToken(headers)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetBearerToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotToken != tt.expectedToken {
				t.Errorf("GetBearerToken() gotToken = %v, want %v", gotToken, tt.expectedToken)
			}
		})
	}
}
