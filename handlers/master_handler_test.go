package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/danilomarques1/fidusserver/database"
	"github.com/danilomarques1/fidusserver/dtos"
	"github.com/danilomarques1/fidusserver/response"
)

const BaseUrl = "http://localhost:8080/fidus"

func dropData(t *testing.T) {
	db := database.Database()
	if _, err := db.Exec(`truncate table fidus_master`); err != nil {
		t.Fatal(err)
	}
}

func TestRegisterService(t *testing.T) {
	defer dropData(t)
	input := bytes.NewReader([]byte(`{"name": "Mocked name", "email":"mock@gmail.com", "password":"thisisasecretpassword"}`))
	req, err := http.NewRequest(http.MethodPost, BaseUrl+"/master/register", input)
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
	req, err := http.NewRequest(http.MethodPost, BaseUrl+"/master/register", input)
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
	input := bytes.NewReader([]byte(`{"name": "Mocked name", "email":"mock@gmail.com", "password":"thisisasecretpassword"}`))
	req, _ := http.NewRequest(http.MethodPost, BaseUrl+"/master/register", input)
	resp, _ := http.DefaultClient.Do(req)
	resp.Body.Close()

	input = bytes.NewReader([]byte(`{"name": "Mocked name", "email":"mock@gmail.com", "password":"thisisasecretpassword"}`))
	req, err := http.NewRequest(http.MethodPost, BaseUrl+"/master/register", input)
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

	respBody := &response.ErrorResponseDto{}
	if err := json.Unmarshal(b, respBody); err != nil {
		t.Fatal(err)
	}
	if respBody.Message != "Email already taken" {
		t.Fatalf("Wrong message returned: %v\n", respBody.Message)
	}
}
