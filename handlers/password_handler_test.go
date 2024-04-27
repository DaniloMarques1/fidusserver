package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/danilomarques1/fidusserver/dtos"
)

func TestStorePassword(t *testing.T) {
	defer dropData(t)
	// create master
	input := `{"name": "Mocked name", "email":"mock@gmail.com", "password":"thisisasecretpassword"}`
	req, _ := http.NewRequest(http.MethodPost, baseUrl+"/master/register", bytes.NewReader([]byte(input)))
	resp, _ := http.DefaultClient.Do(req)
	resp.Body.Close()

	// auth master
	input = `{"email": "mock@gmail.com", "password":"thisisasecretpassword"}`
	req, _ = http.NewRequest(http.MethodPost, baseUrl+"/master/authenticate", bytes.NewReader([]byte(input)))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	authResponse := &dtos.AuthenticateResponseDto{}
	json.Unmarshal(b, authResponse)

	input = `{"key": "somekey", "password":"somepassword"}`
	req, err = http.NewRequest(http.MethodPost, baseUrl+"/password/store", bytes.NewReader([]byte(input)))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("Authorization", "Bearer "+authResponse.AccessToken)
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		t.Fatalf("Wrong status code returned: %v\n", resp.StatusCode)
	}

	//TODO: retrieve the key
	req, err = http.NewRequest(http.MethodGet, baseUrl+"/password/retrieve", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("Authorization", "Bearer "+authResponse.AccessToken)
	query := req.URL.Query()
	query.Add("key", "somekey")
	req.URL.RawQuery = query.Encode()
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Wrong status code returned %v\n", resp.StatusCode)
	}
}

func TestRetrievePassword(t *testing.T) {
	defer dropData(t)
	// create master
	input := `{"name": "Mocked name", "email":"mock@gmail.com", "password":"thisisasecretpassword"}`
	req, _ := http.NewRequest(http.MethodPost, baseUrl+"/master/register", bytes.NewReader([]byte(input)))
	resp, _ := http.DefaultClient.Do(req)
	resp.Body.Close()

	// auth master
	input = `{"email": "mock@gmail.com", "password":"thisisasecretpassword"}`
	req, _ = http.NewRequest(http.MethodPost, baseUrl+"/master/authenticate", bytes.NewReader([]byte(input)))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	authResponse := &dtos.AuthenticateResponseDto{}
	json.Unmarshal(b, authResponse)

	input = `{"key": "somekey", "password":"somepassword"}`
	req, _ = http.NewRequest(http.MethodPost, baseUrl+"/password/store", bytes.NewReader([]byte(input)))
	req.Header.Add("Authorization", "Bearer "+authResponse.AccessToken)
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		t.Fatalf("Wrong status code returned: %v\n", resp.StatusCode)
	}

	req, err = http.NewRequest(http.MethodGet, baseUrl+"/password/retrieve", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("Authorization", "Bearer "+authResponse.AccessToken)
	query := req.URL.Query()
	query.Add("key", "somekey")
	req.URL.RawQuery = query.Encode()
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Wrong status code returned %v\n", resp.StatusCode)
	}
	b, err = io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	respBody := &dtos.RetrievePasswordResponseDto{}
	json.Unmarshal(b, respBody)
	if respBody.Key != "somekey" {
		t.Fatalf("Wrong key returned %v\n", respBody.Key)
	}
}
