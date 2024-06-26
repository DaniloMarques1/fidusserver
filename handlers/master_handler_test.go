package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/danilomarques1/fidusserver/database"
	"github.com/danilomarques1/fidusserver/dtos"
)

const baseUrl = "http://localhost:8080/fidus"

// it cleans the database and should be called
// after running each test
func dropData(t *testing.T) {
	t.Setenv("DATABASE_URI", "postgresql://fitz:fitz@localhost:5432/fidus?sslmode=disable")
	db := database.Database()
	if _, err := db.Exec(`truncate table fidus_master cascade; truncate table fidus_password;`); err != nil {
		t.Fatal(err)
	}
}

func TestRegisterService(t *testing.T) {
	defer dropData(t)
	input := bytes.NewReader([]byte(`{"name": "Mocked name", "email":"mock@gmail.com", "password":"Mock@@123"}`))
	req, err := http.NewRequest(http.MethodPost, baseUrl+"/master/register", input)
	if err != nil {
		t.Fatal(err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("Wrong status code returned. Expected %v got %v\n", http.StatusCreated, resp.StatusCode)
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	respBody := &dtos.CreateMasterResponseDto{}
	if err := json.Unmarshal(b, respBody); err != nil {
		t.Fatal(err)
	}
	if respBody.ID == "" {
		t.Fatal("ID should be defined")
	}

}

func TestRegisterServiceEmptyBody(t *testing.T) {
	defer dropData(t)
	input := bytes.NewReader([]byte(``))
	req, err := http.NewRequest(http.MethodPost, baseUrl+"/master/register", input)
	if err != nil {
		t.Fatal(err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("Wrong status code returned. Expected %v got %v\n", http.StatusBadRequest, resp.StatusCode)
	}
}

func TestRegisterServiceInvalidEmail(t *testing.T) {
	defer dropData(t)
	input := `{"name": "Mocked name", "email":"mock", "password":"Mock@@123"}`
	req, err := http.NewRequest(http.MethodPost, baseUrl+"/master/register", bytes.NewReader([]byte(input)))
	if err != nil {
		t.Fatal(err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("Wrong status code returned. Expected %v got %v\n", http.StatusBadRequest, resp.StatusCode)
	}
}

func TestRegisterServiceEmptyPassword(t *testing.T) {
	defer dropData(t)
	input := `{"name": "Mocked name", "email":"mock@mail.com", "password":""}`
	req, err := http.NewRequest(http.MethodPost, baseUrl+"/master/register", bytes.NewReader([]byte(input)))
	if err != nil {
		t.Fatal(err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("Wrong status code returned. Expected %v got %v\n", http.StatusBadRequest, resp.StatusCode)
	}
}

func TestRegisterServiceInvalidPassword(t *testing.T) {
	defer dropData(t)
	input := `{"name": "Mocked name", "email":"mock@mail.com", "password":"1234"}`
	req, err := http.NewRequest(http.MethodPost, baseUrl+"/master/register", bytes.NewReader([]byte(input)))
	if err != nil {
		t.Fatal(err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("Wrong status code returned. Expected %v got %v\n", http.StatusBadRequest, resp.StatusCode)
	}
}

func TestRegisterServiceEmailAlreadyTaken(t *testing.T) {
	defer dropData(t)
	input := `{"name": "Mocked name", "email":"mock@gmail.com", "password":"Mock@@123"}`
	req, _ := http.NewRequest(http.MethodPost, baseUrl+"/master/register", bytes.NewReader([]byte(input)))
	resp, _ := http.DefaultClient.Do(req)
	if resp.StatusCode != http.StatusCreated {
		t.Fatal("Wrong status code when creating the master")
	}
	resp.Body.Close()

	input = `{"name": "Mocked name", "email":"mock@gmail.com", "password":"Mock@@123"}`
	req, err := http.NewRequest(http.MethodPost, baseUrl+"/master/register", bytes.NewReader([]byte(input)))
	if err != nil {
		t.Fatal(err)
	}
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("Wrong status code returned. Expected %v got %v\n", http.StatusCreated, resp.StatusCode)
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	respBody := &dtos.ErrorResponseDto{}
	if err := json.Unmarshal(b, respBody); err != nil {
		t.Fatal(err)
	}
	if respBody.Message != "Email already taken" {
		t.Fatalf("Wrong message returned: %v\n", respBody.Message)
	}
}

func TestMasterAuthenticate(t *testing.T) {
	defer dropData(t)
	// creating a master
	input := `{"name": "Mocked name", "email":"mock@gmail.com", "password":"Mock@@123"}`
	req, _ := http.NewRequest(http.MethodPost, baseUrl+"/master/register", bytes.NewReader([]byte(input)))
	resp, _ := http.DefaultClient.Do(req)
	if resp.StatusCode != http.StatusCreated {
		t.Fatal("Wrong status code when creating the master")
	}
	resp.Body.Close()

	input = `{"email": "mock@gmail.com", "password":"Mock@@123"}`
	req, err := http.NewRequest(http.MethodPost, baseUrl+"/master/authenticate", bytes.NewReader([]byte(input)))
	if err != nil {
		t.Fatal(err)
	}
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Wrong status code returned. Expected: %v got: %v\n", http.StatusOK, resp.StatusCode)
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	respBody := &dtos.AuthenticateResponseDto{}
	if err := json.Unmarshal(b, respBody); err != nil {
		t.Fatal(err)
	}
	if len(respBody.AccessToken) == 0 {
		t.Fatalf("Access token not returned")
	}
	if respBody.ExpiresAt <= 0 {
		t.Fatalf("Access token not returned")
	}
}

func TestMasterAuthenticateInvalidEmail(t *testing.T) {
	defer dropData(t)
	// creating a master
	input := `{"name": "Mocked name", "email":"mock@gmail.com", "password":"Mock@@123"}`
	req, _ := http.NewRequest(http.MethodPost, baseUrl+"/master/register", bytes.NewReader([]byte(input)))
	resp, _ := http.DefaultClient.Do(req)
	if resp.StatusCode != http.StatusCreated {
		t.Fatal("Wrong status code returned when creating the master")
	}
	resp.Body.Close()

	input = `{"email": "mockcom", "password":"Mock@@123"}`
	req, err := http.NewRequest(http.MethodPost, baseUrl+"/master/authenticate", bytes.NewReader([]byte(input)))
	if err != nil {
		t.Fatal(err)
	}
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("Wrong status code returned. Expected: %v got: %v\n", http.StatusOK, resp.StatusCode)
	}
}

func TestMasterAuthenticateWrongEmail(t *testing.T) {
	defer dropData(t)
	input := `{"email": "mock@gmail.com", "password":"Mock@@123"}`
	req, err := http.NewRequest(http.MethodPost, baseUrl+"/master/authenticate", bytes.NewReader([]byte(input)))
	if err != nil {
		t.Fatal(err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("Wrong status code returned. Expected: %v got: %v\n", http.StatusBadRequest, resp.StatusCode)
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	errorDto := &dtos.ErrorResponseDto{}
	if err := json.Unmarshal(b, errorDto); err != nil {
		t.Fatal(err)
	}
	if errorDto.Message != "Incorrect credentials" {
		t.Fatal("Wrong message returned")
	}
}

func TestMasterAuthenticateWrongPassword(t *testing.T) {
	defer dropData(t)
	// creating a master
	input := `{"name": "Mocked name", "email":"mock@gmail.com", "password":"Mock@@123"}`
	req, _ := http.NewRequest(http.MethodPost, baseUrl+"/master/register", bytes.NewReader([]byte(input)))
	resp, _ := http.DefaultClient.Do(req)
	if resp.StatusCode != http.StatusCreated {
		t.Fatal("Wrong status code when creating the master")
	}
	resp.Body.Close()

	input = `{"email": "mock@gmail.com", "password":"Mock@@124"}`
	req, err := http.NewRequest(http.MethodPost, baseUrl+"/master/authenticate", bytes.NewReader([]byte(input)))
	if err != nil {
		t.Fatal(err)
	}
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("Wrong status code returned. Expected: %v got: %v\n", http.StatusBadRequest, resp.StatusCode)
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	errorDto := &dtos.ErrorResponseDto{}
	if err := json.Unmarshal(b, errorDto); err != nil {
		t.Fatal(err)
	}
	if errorDto.Message != "Incorrect credentials" {
		t.Fatalf("Wrong message returned %v", errorDto.Message)
	}
}

func TestResetMasterPassword(t *testing.T) {
	defer dropData(t)
	// creating a master
	input := `{"name": "Mocked name", "email":"mock@gmail.com", "password":"Mock@@123"}`
	req, _ := http.NewRequest(http.MethodPost, baseUrl+"/master/register", bytes.NewReader([]byte(input)))
	resp, _ := http.DefaultClient.Do(req)
	resp.Body.Close()

	// authenticating
	input = `{"email": "mock@gmail.com", "password":"Mock@@123"}`
	http.NewRequest(http.MethodPost, baseUrl+"/master/authenticate", bytes.NewReader([]byte(input)))
	http.DefaultClient.Do(req)
	resp.Body.Close()

	// update
	input = `{"email": "mock@gmail.com", "old_password":"Mock@@123", "new_password": "Mock@@124"}`
	req, _ = http.NewRequest(http.MethodPut, baseUrl+"/master/reset/password", bytes.NewReader([]byte(input)))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if resp.StatusCode != http.StatusNoContent {
		t.Fatalf("Wrong status code returned. expected: %v got: %v\n", http.StatusNoContent, resp.StatusCode)
	}

	// authenticating
	input = `{"email": "mock@gmail.com", "password":"Mock@@124"}`
	req, err = http.NewRequest(http.MethodPost, baseUrl+"/master/authenticate", bytes.NewReader([]byte(input)))
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Wrong status code returned. expected: %v got: %v\n", http.StatusOK, resp.StatusCode)
	}
}

func TestResetMasterPasswordWrongOldPassword(t *testing.T) {
	defer dropData(t)
	// creating a master
	input := `{"name": "Mocked name", "email":"mock@gmail.com", "password":"Mock@@123"}`
	req, _ := http.NewRequest(http.MethodPost, baseUrl+"/master/register", bytes.NewReader([]byte(input)))
	resp, _ := http.DefaultClient.Do(req)
	resp.Body.Close()

	// authenticating
	input = `{"email": "mock@gmail.com", "password":"Mock@@123"}`
	http.NewRequest(http.MethodPost, baseUrl+"/master/authenticate", bytes.NewReader([]byte(input)))
	http.DefaultClient.Do(req)
	resp.Body.Close()

	// update
	input = `{"email": "mock@gmail.com", "old_password":"Mock@@122", "new_password": "Mock@@124"}`
	req, _ = http.NewRequest(http.MethodPut, baseUrl+"/master/reset/password", bytes.NewReader([]byte(input)))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("Wrong status code returned. expected: %v got: %v\n", http.StatusNoContent, resp.StatusCode)
	}

	// authenticating
	input = `{"email": "mock@gmail.com", "password":"Mock@@123"}`
	req, err = http.NewRequest(http.MethodPost, baseUrl+"/master/authenticate", bytes.NewReader([]byte(input)))
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Wrong status code returned. expected: %v got: %v\n", http.StatusOK, resp.StatusCode)
	}
}
