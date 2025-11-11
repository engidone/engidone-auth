package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type SigninRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type SigninResponse struct {
	Success   bool   `json:"success"`
	Message   string `json:"message"`
	UserID    string `json:"user_id,omitempty"`
	Username  string `json:"username,omitempty"`
	Email     string `json:"email,omitempty"`
	Token     string `json:"token,omitempty"`
	ExpiresAt int64  `json:"expires_at,omitempty"`
}

type ValidateTokenRequest struct {
	Token string `json:"token"`
}

type ValidateTokenResponse struct {
	Valid    bool   `json:"valid"`
	Message  string `json:"message"`
	UserID   string `json:"user_id,omitempty"`
	Username string `json:"username,omitempty"`
	Email    string `json:"email,omitempty"`
}

type GetUserRequest struct {
	UserID string `json:"user_id"`
}

type GetUserResponse struct {
	Success   bool   `json:"success"`
	Message   string `json:"message"`
	UserID    string `json:"user_id,omitempty"`
	Username  string `json:"username,omitempty"`
	Email     string `json:"email,omitempty"`
	CreatedAt int64  `json:"created_at,omitempty"`
	UpdatedAt int64  `json:"updated_at,omitempty"`
}

func main() {
	fmt.Println("=== Signin Service Client (HTTP) ===")
	fmt.Println("Probando servidor unificado en :8080")
	fmt.Println()

	baseURL := "http://localhost:8080"

	// Test 1: Signin exitoso
	fmt.Println("1. Probando Signin exitoso...")
	signinReq := SigninRequest{
		Username: "admin",
		Password: "password123",
	}

	jsonBody, _ := json.Marshal(signinReq)
	resp, err := http.Post(baseURL+"/signin", "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		fmt.Printf("   Error en signin: %v\n", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("   Error leyendo respuesta: %v\n", err)
		return
	}

	var signinRsp SigninResponse
	if err := json.Unmarshal(body, &signinRsp); err != nil {
		fmt.Printf("   Error parseando respuesta: %v\n", err)
		fmt.Printf("   Raw response: %s\n", string(body))
	} else {
		if signinRsp.Success {
			fmt.Printf("   ✅ Signin exitoso!\n")
			fmt.Printf("   UserID: %s\n", signinRsp.UserID)
			fmt.Printf("   Username: %s\n", signinRsp.Username)
			fmt.Printf("   Email: %s\n", signinRsp.Email)
			fmt.Printf("   Token: %s\n", signinRsp.Token)
			fmt.Printf("   ExpiresAt: %d\n", signinRsp.ExpiresAt)
		} else {
			fmt.Printf("   ❌ Signin fallido: %s\n", signinRsp.Message)
		}
	}
	fmt.Println()

	// Guardar token para pruebas siguientes
	validToken := signinRsp.Token
	validUserID := signinRsp.UserID

	// Test 2: Signin fallido
	fmt.Println("2. Probando Signin con credenciales incorrectas...")
	invalidReq := SigninRequest{
		Username: "admin",
		Password: "wrongpassword",
	}

	jsonBody, _ = json.Marshal(invalidReq)
	resp, err = http.Post(baseURL+"/signin", "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		fmt.Printf("   Error en signin: %v\n", err)
	} else {
		defer resp.Body.Close()
		body, _ = io.ReadAll(resp.Body)

		var invalidRsp SigninResponse
		if err := json.Unmarshal(body, &invalidRsp); err == nil {
			if !invalidRsp.Success {
				fmt.Printf("   ✅ Signin correctamente rechazado: %s\n", invalidRsp.Message)
			} else {
				fmt.Printf("   ❌ No debería haber funcionado!\n")
			}
		}
	}
	fmt.Println()

	// Test 3: Validar token
	if validToken != "" {
		fmt.Println("3. Probando ValidateToken...")
		validateReq := ValidateTokenRequest{
			Token: validToken,
		}

		jsonBody, _ = json.Marshal(validateReq)
		resp, err := http.Post(baseURL+"/validate-token", "application/json", bytes.NewBuffer(jsonBody))
		if err != nil {
			fmt.Printf("   Error validando token: %v\n", err)
		} else {
			defer resp.Body.Close()
			body, _ = io.ReadAll(resp.Body)

			var validateRsp ValidateTokenResponse
			if err := json.Unmarshal(body, &validateRsp); err == nil {
				if validateRsp.Valid {
					fmt.Printf("   ✅ Token válido!\n")
					fmt.Printf("   UserID: %s\n", validateRsp.UserID)
					fmt.Printf("   Username: %s\n", validateRsp.Username)
					fmt.Printf("   Email: %s\n", validateRsp.Email)
				} else {
					fmt.Printf("   ❌ Token inválido: %s\n", validateRsp.Message)
				}
			}
		}
		fmt.Println()
	}

	// Test 4: GetUser
	if validUserID != "" {
		fmt.Println("4. Probando GetUser...")
		getUserReq := GetUserRequest{
			UserID: validUserID,
		}

		jsonBody, _ = json.Marshal(getUserReq)
		resp, err := http.Post(baseURL+"/get-user", "application/json", bytes.NewBuffer(jsonBody))
		if err != nil {
			fmt.Printf("   Error obteniendo usuario: %v\n", err)
		} else {
			defer resp.Body.Close()
			body, _ = io.ReadAll(resp.Body)

			var getUserRsp GetUserResponse
			if err := json.Unmarshal(body, &getUserRsp); err == nil {
				if getUserRsp.Success {
					fmt.Printf("   ✅ Usuario obtenido!\n")
					fmt.Printf("   UserID: %s\n", getUserRsp.UserID)
					fmt.Printf("   Username: %s\n", getUserRsp.Username)
					fmt.Printf("   Email: %s\n", getUserRsp.Email)
					fmt.Printf("   CreatedAt: %d\n", getUserRsp.CreatedAt)
					fmt.Printf("   UpdatedAt: %d\n", getUserRsp.UpdatedAt)
				} else {
					fmt.Printf("   ❌ Error obteniendo usuario: %s\n", getUserRsp.Message)
				}
			}
		}
		fmt.Println()
	}

	// Test 5: Múltiples usuarios
	fmt.Println("5. Probando signin con diferentes usuarios...")
	testUsers := []struct {
		username string
		password string
	}{
		{"testuser", "test123"},
		{"john", "john123"},
		{"nonexistent", "wrong"},
	}

	for _, user := range testUsers {
		req := SigninRequest{
			Username: user.username,
			Password: user.password,
		}

		jsonBody, _ = json.Marshal(req)
		resp, err := http.Post(baseURL+"/signin", "application/json", bytes.NewBuffer(jsonBody))
		if err != nil {
			fmt.Printf("   %s/%s: Error - %v\n", user.username, user.password, err)
			continue
		}

		defer resp.Body.Close()
		body, _ = io.ReadAll(resp.Body)

		var rsp SigninResponse
		if err := json.Unmarshal(body, &rsp); err == nil {
			if rsp.Success {
				fmt.Printf("   ✅ %s/%s: Success (ID: %s)\n", user.username, user.password, rsp.UserID)
			} else {
				fmt.Printf("   ❌ %s/%s: %s\n", user.username, user.password, rsp.Message)
			}
		}

		time.Sleep(500 * time.Millisecond)
	}
	fmt.Println()

	fmt.Println("=== Pruebas completadas ===")
}