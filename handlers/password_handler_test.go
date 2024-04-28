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

	// auth master
	input = `{"email": "mock@gmail.com", "password":"thisisasecretpassword"}`
	req, _ = http.NewRequest(http.MethodPost, baseUrl+"/master/authenticate", bytes.NewReader([]byte(input)))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Error(err)
	}
	b, err := io.ReadAll(resp.Body)
	authResponse := &dtos.AuthenticateResponseDto{}
	json.Unmarshal(b, authResponse)

	input = `{"key": "somekey", "password":"somepassword"}`
	req, err = http.NewRequest(http.MethodPost, baseUrl+"/password/store", bytes.NewReader([]byte(input)))
	if err != nil {
		t.Error(err)
	}
	req.Header.Add("Authorization", "Bearer "+authResponse.AccessToken)
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Error(err)
	}

	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("Wrong status code returned: %v\n", resp.StatusCode)
	}

	req, err = http.NewRequest(http.MethodGet, baseUrl+"/password/retrieve", nil)
	if err != nil {
		t.Error(err)
	}
	req.Header.Add("Authorization", "Bearer "+authResponse.AccessToken)
	query := req.URL.Query()
	query.Add("key", "somekey")
	req.URL.RawQuery = query.Encode()
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Error(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Wrong status code returned %v\n", resp.StatusCode)
	}
}

func TestStorePasswordWrongToken(t *testing.T) {
	defer dropData(t)
	accessToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IlF1aW5jeSBMYXJzb24iLCJpYXQiOjE1MTYyMzkwMjJ9.WcPGXClpKD7Bc1C0CCDA1060E2GGlTfamrd8 - W0ghBE"
	// create master
	input := `{"name": "Mocked name", "email":"mock@gmail.com", "password":"thisisasecretpassword"}`
	req, _ := http.NewRequest(http.MethodPost, baseUrl+"/master/register", bytes.NewReader([]byte(input)))
	resp, _ := http.DefaultClient.Do(req)

	// auth master
	input = `{"email": "mock@gmail.com", "password":"thisisasecretpassword"}`
	req, _ = http.NewRequest(http.MethodPost, baseUrl+"/master/authenticate", bytes.NewReader([]byte(input)))
	http.DefaultClient.Do(req)

	// create a password
	input = `{"key": "somekey", "password":"somepassword"}`
	req, err := http.NewRequest(http.MethodPost, baseUrl+"/password/store", bytes.NewReader([]byte(input)))
	if err != nil {
		t.Error(err)
	}
	req.Header.Add("Authorization", "Bearer "+accessToken)
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Error(err)
	}

	if resp.StatusCode != http.StatusForbidden {
		t.Errorf("Wrong status code returned %v\n", resp.StatusCode)
	}
}

func TestStorePasswordExpiredToken(t *testing.T) {
	defer dropData(t)
	expiredToken := "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJtYXN0ZXJfaWQiOiI5NjI5ODhjYS0yMmQ5LTQ2ZDktYmNiOC1iYWI5NDllN2UyNzUiLCJtYXN0ZXJfZW1haWwiOiJtb2NrQGdtYWlsLmNvbSIsImV4cCI6MTcxNDIzMjgwNX0.xRY-TxkWitEHAj8Ow5i308d3iE_yQoy7JAK4wJwToNXLXwORs3A1QpcnUjX8ZiTg05BSo7Hkl7eJxLTRdliDWw"
	// create master
	input := `{"name": "Mocked name", "email":"mock@gmail.com", "password":"thisisasecretpassword"}`
	req, _ := http.NewRequest(http.MethodPost, baseUrl+"/master/register", bytes.NewReader([]byte(input)))
	resp, _ := http.DefaultClient.Do(req)

	// auth master
	input = `{"email": "mock@gmail.com", "password":"thisisasecretpassword"}`
	req, _ = http.NewRequest(http.MethodPost, baseUrl+"/master/authenticate", bytes.NewReader([]byte(input)))
	http.DefaultClient.Do(req)

	// create a password
	input = `{"key": "somekey", "password":"somepassword"}`
	req, err := http.NewRequest(http.MethodPost, baseUrl+"/password/store", bytes.NewReader([]byte(input)))
	if err != nil {
		t.Error(err)
	}
	req.Header.Add("Authorization", "Bearer "+expiredToken)
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Error(err)
	}

	if resp.StatusCode != http.StatusForbidden {
		t.Errorf("Wrong status code returned %v\n", resp.StatusCode)
	}
}

func TestRetrievePassword(t *testing.T) {
	defer dropData(t)
	// create master
	input := `{"name": "Mocked name", "email":"mock@gmail.com", "password":"thisisasecretpassword"}`
	req, _ := http.NewRequest(http.MethodPost, baseUrl+"/master/register", bytes.NewReader([]byte(input)))
	resp, _ := http.DefaultClient.Do(req)

	// auth master
	input = `{"email": "mock@gmail.com", "password":"thisisasecretpassword"}`
	req, _ = http.NewRequest(http.MethodPost, baseUrl+"/master/authenticate", bytes.NewReader([]byte(input)))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Error(err)
	}
	b, err := io.ReadAll(resp.Body)
	authResponse := &dtos.AuthenticateResponseDto{}
	json.Unmarshal(b, authResponse)

	input = `{"key": "somekey", "password":"somepassword"}`
	req, _ = http.NewRequest(http.MethodPost, baseUrl+"/password/store", bytes.NewReader([]byte(input)))
	req.Header.Add("Authorization", "Bearer "+authResponse.AccessToken)
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Error(err)
	}

	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("Wrong status code returned: %v\n", resp.StatusCode)
	}

	req, err = http.NewRequest(http.MethodGet, baseUrl+"/password/retrieve", nil)
	if err != nil {
		t.Error(err)
	}
	req.Header.Add("Authorization", "Bearer "+authResponse.AccessToken)
	query := req.URL.Query()
	query.Add("key", "somekey")
	req.URL.RawQuery = query.Encode()
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Error(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Wrong status code returned %v\n", resp.StatusCode)
	}
	b, err = io.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	respBody := &dtos.RetrievePasswordResponseDto{}
	json.Unmarshal(b, respBody)
	if respBody.Key != "somekey" {
		t.Errorf("Wrong key returned %v\n", respBody.Key)
	}
}

func TestRetrievePasswordWrongKey(t *testing.T) {
	defer dropData(t)
	// create master
	input := `{"name": "Mocked name", "email":"mock@gmail.com", "password":"thisisasecretpassword"}`
	req, _ := http.NewRequest(http.MethodPost, baseUrl+"/master/register", bytes.NewReader([]byte(input)))
	resp, _ := http.DefaultClient.Do(req)

	// auth master
	input = `{"email": "mock@gmail.com", "password":"thisisasecretpassword"}`
	req, _ = http.NewRequest(http.MethodPost, baseUrl+"/master/authenticate", bytes.NewReader([]byte(input)))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Error(err)
	}
	b, err := io.ReadAll(resp.Body)
	authResponse := &dtos.AuthenticateResponseDto{}
	json.Unmarshal(b, authResponse)

	input = `{"key": "somekey", "password":"somepassword"}`
	req, _ = http.NewRequest(http.MethodPost, baseUrl+"/password/store", bytes.NewReader([]byte(input)))
	req.Header.Add("Authorization", "Bearer "+authResponse.AccessToken)
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Error(err)
	}

	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("Wrong status code returned: %v\n", resp.StatusCode)
	}

	req, err = http.NewRequest(http.MethodGet, baseUrl+"/password/retrieve", nil)
	if err != nil {
		t.Error(err)
	}
	req.Header.Add("Authorization", "Bearer "+authResponse.AccessToken)
	query := req.URL.Query()
	query.Add("key", "wrongkey")
	req.URL.RawQuery = query.Encode()
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Error(err)
	}
	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("Wrong status code returned %v\n", resp.StatusCode)
	}
}

func TestDeletePassword(t *testing.T) {
	// create a new master account
	input := `{"name":"mock", "email": "mock@email.com", "password":"12345678"}`
	http.Post(baseUrl+"/master/register", "application/json", bytes.NewReader([]byte(input)))

	// authenticate a master
	input = `{"email": "mock@email.com", "password":"12345678"}`
	resp, _ := http.Post(baseUrl+"/master/authenticate", "application/json", bytes.NewReader([]byte(input)))
	b, _ := io.ReadAll(resp.Body)
	authResponse := &dtos.AuthenticateResponseDto{}
	json.Unmarshal(b, authResponse)

	// create a new password
	input = `{"key": "somekey", "password":"somepassword"}`
	req, _ := http.NewRequest(http.MethodPost, baseUrl+"/password/store", bytes.NewReader([]byte(input)))
	req.Header.Add("Authorization", "Bearer "+authResponse.AccessToken)
	resp, _ = http.DefaultClient.Do(req)

	// delete a password
	req, _ = http.NewRequest(http.MethodDelete, baseUrl+"/password/delete", nil)
	query := req.URL.Query()
	query.Add("key", "somekey")
	req.URL.RawQuery = query.Encode()
	req.Header.Add("Authorization", "Bearer "+authResponse.AccessToken)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Error(err)
	}
	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("Wrong status code returned: %v\n", resp.StatusCode)
	}

	// search for the password
	req, _ = http.NewRequest(http.MethodGet, baseUrl+"/password/retrieve", nil)
	query = req.URL.Query()
	query.Add("key", "somekey")
	req.URL.RawQuery = query.Encode()
	req.Header.Add("Authorization", "Bearer "+authResponse.AccessToken)
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Error(err)
	}
	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("Wrong status code returned: %v\n", resp.StatusCode)
	}
}

func TestUpdatePassword(t *testing.T) {
	defer dropData(t)
	// create a new master account
	input := `{"name":"mock", "email": "mock@email.com", "password":"12345678"}`
	http.Post(baseUrl+"/master/register", "application/json", bytes.NewReader([]byte(input)))

	// authenticate a master
	input = `{"email": "mock@email.com", "password":"12345678"}`
	resp, _ := http.Post(baseUrl+"/master/authenticate", "application/json", bytes.NewReader([]byte(input)))
	b, _ := io.ReadAll(resp.Body)
	authResponse := &dtos.AuthenticateResponseDto{}
	json.Unmarshal(b, authResponse)

	// create a new password
	input = `{"key": "somekey", "password":"somepassword"}`
	req, _ := http.NewRequest(http.MethodPost, baseUrl+"/password/store", bytes.NewReader([]byte(input)))
	req.Header.Add("Authorization", "Bearer "+authResponse.AccessToken)
	http.DefaultClient.Do(req)

	// update password
	input = `{"password":"updatepassword"}`
	req, _ = http.NewRequest(http.MethodPut, baseUrl+"/password/update", bytes.NewReader([]byte(input)))
	query := req.URL.Query()
	query.Add("key", "somekey")
	req.Header.Add("Authorization", "Bearer "+authResponse.AccessToken)
	req.URL.RawQuery = query.Encode()
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Error(err)
	}
	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("Wrong status code returned %v\n", resp.StatusCode)
	}

	// search for the password
	req, _ = http.NewRequest(http.MethodGet, baseUrl+"/password/retrieve", nil)
	query = req.URL.Query()
	query.Add("key", "somekey")
	req.URL.RawQuery = query.Encode()
	req.Header.Add("Authorization", "Bearer "+authResponse.AccessToken)
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Error(err)
	}
	b, err = io.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Wrong status code returned: %v\n", resp.StatusCode)
	}

	retrieveResponse := &dtos.RetrievePasswordResponseDto{}
	json.Unmarshal(b, retrieveResponse)
	if retrieveResponse.Password != "updatepassword" {
		t.Errorf("Wrong password value returned: %v\n", retrieveResponse.Password)
	}
}
