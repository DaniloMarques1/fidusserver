package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/danilomarques1/fidusserver/dtos"
)

func createAndAuthenticateMaster() (string, error) {
	// create master
	input := `{"name": "Mocked name", "email":"mock@gmail.com", "password":"thisisasecretpassword"}`
	req, err := http.NewRequest(http.MethodPost, baseUrl+"/master/register", bytes.NewReader([]byte(input)))
	if err != nil {
		return "", err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", nil
	}

	// auth master
	input = `{"email": "mock@gmail.com", "password":"thisisasecretpassword"}`
	req, err = http.NewRequest(http.MethodPost, baseUrl+"/master/authenticate", bytes.NewReader([]byte(input)))
	if err != nil {
		return "", err
	}
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	authResponse, err := readResponse[dtos.AuthenticateResponseDto](resp)
	if err != nil {
		return "", err
	}
	return authResponse.AccessToken, nil
}

func readResponse[T any](resp *http.Response) (*T, error) {
	var t T
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(b, &t); err != nil {
		return nil, err
	}
	return &t, nil
}

func TestStorePassword(t *testing.T) {
	defer dropData(t)
	accessToken, err := createAndAuthenticateMaster()
	if err != nil {
		t.Fatal(err)
	}

	// creating a password
	input := `{"key": "somekey", "password":"somepassword"}`
	req, err := http.NewRequest(http.MethodPost, baseUrl+"/password/store", bytes.NewReader([]byte(input)))
	if err != nil {
		t.Error(err)
	}
	req.Header.Add("Authorization", "Bearer "+accessToken)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Error(err)
	}

	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("Wrong status code returned: %v\n", resp.StatusCode)
	}

	// tries to retrieve the created password
	req, err = http.NewRequest(http.MethodGet, baseUrl+"/password/retrieve", nil)
	if err != nil {
		t.Error(err)
	}
	req.Header.Add("Authorization", "Bearer "+accessToken)
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

func TestStorePasswordEmptyKey(t *testing.T) {
	defer dropData(t)
	accessToken, err := createAndAuthenticateMaster()
	if err != nil {
		t.Fatal(err)
	}

	input := `{"key": "", "password":"somepassword"}`
	req, err := http.NewRequest(http.MethodPost, baseUrl+"/password/store", bytes.NewReader([]byte(input)))
	if err != nil {
		t.Error(err)
	}
	req.Header.Add("Authorization", "Bearer "+accessToken)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Error(err)
	}

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Wrong status code returned: %v\n", resp.StatusCode)
	}
}

func TestStorePasswordWrongToken(t *testing.T) {
	defer dropData(t)
	wrongAccessToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IlF1aW5jeSBMYXJzb24iLCJpYXQiOjE1MTYyMzkwMjJ9.WcPGXClpKD7Bc1C0CCDA1060E2GGlTfamrd8 - W0ghBE"
	if _, err := createAndAuthenticateMaster(); err != nil {
		t.Fatal(err)
	}
	// create a password
	input := `{"key": "somekey", "password":"somepassword"}`
	req, err := http.NewRequest(http.MethodPost, baseUrl+"/password/store", bytes.NewReader([]byte(input)))
	if err != nil {
		t.Error(err)
	}
	req.Header.Add("Authorization", "Bearer "+wrongAccessToken)
	resp, err := http.DefaultClient.Do(req)
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
	_, err := createAndAuthenticateMaster()
	if err != nil {
		t.Fatal(err)
	}

	// create a password
	input := `{"key": "somekey", "password":"somepassword"}`
	req, err := http.NewRequest(http.MethodPost, baseUrl+"/password/store", bytes.NewReader([]byte(input)))
	if err != nil {
		t.Error(err)
	}
	req.Header.Add("Authorization", "Bearer "+expiredToken)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Error(err)
	}

	if resp.StatusCode != http.StatusForbidden {
		t.Errorf("Wrong status code returned %v\n", resp.StatusCode)
	}
}

func TestStorePasswordNoToken(t *testing.T) {
	defer dropData(t)
	// create master
	_, err := createAndAuthenticateMaster()
	if err != nil {
		t.Fatal(err)
	}

	// create a password
	input := `{"key": "somekey", "password":"somepassword"}`
	resp, err := http.Post(baseUrl+"/password/store", "application/json", bytes.NewReader([]byte(input)))
	if err != nil {
		t.Error(err)
	}

	if resp.StatusCode != http.StatusForbidden {
		t.Errorf("Wrong status code returned %v\n", resp.StatusCode)
	}
}

func TestStorePasswordEmptyAuthorizationToken(t *testing.T) {
	defer dropData(t)
	_, err := createAndAuthenticateMaster()
	if err != nil {
		t.Fatal(err)
	}

	// create a password
	input := `{"key": "somekey", "password":"somepassword"}`
	req, err := http.NewRequest(http.MethodPost, baseUrl+"/password/store", bytes.NewReader([]byte(input)))
	if err != nil {
		t.Error(err)
	}
	req.Header.Add("Authorization", "")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Error(err)
	}

	if resp.StatusCode != http.StatusForbidden {
		t.Errorf("Wrong status code returned %v\n", resp.StatusCode)
	}
}

func TestStorePasswordEmptyBearerToken(t *testing.T) {
	defer dropData(t)
	_, err := createAndAuthenticateMaster()
	if err != nil {
		t.Fatal(err)
	}

	// create a password
	input := `{"key": "somekey", "password":"somepassword"}`
	req, err := http.NewRequest(http.MethodPost, baseUrl+"/password/store", bytes.NewReader([]byte(input)))
	if err != nil {
		t.Error(err)
	}
	req.Header.Add("Authorization", "Bearer")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Error(err)
	}

	if resp.StatusCode != http.StatusForbidden {
		t.Errorf("Wrong status code returned %v\n", resp.StatusCode)
	}
}

func TestStorePasswordKeyAlreadyUsed(t *testing.T) {
	defer dropData(t)
	accessToken, err := createAndAuthenticateMaster()
	if err != nil {
		t.Fatal(err)
	}

	// create a password
	input := `{"key": "somekey", "password":"somepassword"}`
	req, _ := http.NewRequest(http.MethodPost, baseUrl+"/password/store", bytes.NewReader([]byte(input)))
	req.Header.Add("Authorization", "Bearer "+accessToken)
	http.DefaultClient.Do(req)

	// create a new  password with the same key
	input = `{"key": "somekey", "password":"anotherpassword"}`
	req, err = http.NewRequest(http.MethodPost, baseUrl+"/password/store", bytes.NewReader([]byte(input)))
	if err != nil {
		t.Error(err)
	}
	req.Header.Add("Authorization", "Bearer "+accessToken)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Error(err)
	}

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Wrong status code returned %v\n", resp.StatusCode)
	}
}

func TestRetrievePassword(t *testing.T) {
	defer dropData(t)
	accessToken, err := createAndAuthenticateMaster()
	if err != nil {
		t.Fatal(err)
	}

	input := `{"key": "somekey", "password":"somepassword"}`
	req, _ := http.NewRequest(http.MethodPost, baseUrl+"/password/store", bytes.NewReader([]byte(input)))
	req.Header.Add("Authorization", "Bearer "+accessToken)
	resp, err := http.DefaultClient.Do(req)
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
	req.Header.Add("Authorization", "Bearer "+accessToken)
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
	respBody, err := readResponse[dtos.RetrievePasswordResponseDto](resp)
	if err != nil {
		t.Error(err)
	}
	if respBody.Key != "somekey" {
		t.Errorf("Wrong key returned %v\n", respBody.Key)
	}
}

func TestRetrievePasswordWrongKey(t *testing.T) {
	defer dropData(t)
	accessToken, err := createAndAuthenticateMaster()
	if err != nil {
		t.Fatal(err)
	}

	input := `{"key": "somekey", "password":"somepassword"}`
	req, _ := http.NewRequest(http.MethodPost, baseUrl+"/password/store", bytes.NewReader([]byte(input)))
	req.Header.Add("Authorization", "Bearer "+accessToken)
	resp, err := http.DefaultClient.Do(req)
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
	req.Header.Add("Authorization", "Bearer "+accessToken)
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
	defer dropData(t)
	accessToken, err := createAndAuthenticateMaster()
	if err != nil {
		t.Fatal(err)
	}

	// create a new password
	input := `{"key": "somekey", "password":"somepassword"}`
	req, _ := http.NewRequest(http.MethodPost, baseUrl+"/password/store", bytes.NewReader([]byte(input)))
	req.Header.Add("Authorization", "Bearer "+accessToken)
	resp, _ := http.DefaultClient.Do(req)

	// delete a password
	req, _ = http.NewRequest(http.MethodDelete, baseUrl+"/password/delete", nil)
	query := req.URL.Query()
	query.Add("key", "somekey")
	req.URL.RawQuery = query.Encode()
	req.Header.Add("Authorization", "Bearer "+accessToken)
	resp, err = http.DefaultClient.Do(req)
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
	req.Header.Add("Authorization", "Bearer "+accessToken)
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
	accessToken, err := createAndAuthenticateMaster()
	if err != nil {
		t.Fatal(err)
	}

	// create a new password
	input := `{"key": "somekey", "password":"somepassword"}`
	req, _ := http.NewRequest(http.MethodPost, baseUrl+"/password/store", bytes.NewReader([]byte(input)))
	req.Header.Add("Authorization", "Bearer "+accessToken)
	http.DefaultClient.Do(req)

	// update password
	input = `{"password":"updatepassword"}`
	req, _ = http.NewRequest(http.MethodPut, baseUrl+"/password/update", bytes.NewReader([]byte(input)))
	query := req.URL.Query()
	query.Add("key", "somekey")
	req.Header.Add("Authorization", "Bearer "+accessToken)
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
	req.Header.Add("Authorization", "Bearer "+accessToken)
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Error(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Wrong status code returned: %v\n", resp.StatusCode)
	}

	retrieveResponse, err := readResponse[dtos.RetrievePasswordResponseDto](resp)
	if err != nil {
		t.Error(err)
	}
	if retrieveResponse.Password != "updatepassword" {
		t.Errorf("Wrong password value returned: %v\n", retrieveResponse.Password)
	}
}

func TestUpdatePasswordEmptyPassword(t *testing.T) {
	defer dropData(t)
	accessToken, err := createAndAuthenticateMaster()
	if err != nil {
		t.Fatal(err)
	}

	// create a new password
	input := `{"key": "somekey", "password":"somepassword"}`
	req, _ := http.NewRequest(http.MethodPost, baseUrl+"/password/store", bytes.NewReader([]byte(input)))
	req.Header.Add("Authorization", "Bearer "+accessToken)
	http.DefaultClient.Do(req)

	// update password
	input = `{"password":""}`
	req, _ = http.NewRequest(http.MethodPut, baseUrl+"/password/update", bytes.NewReader([]byte(input)))
	query := req.URL.Query()
	query.Add("key", "somekey")
	req.Header.Add("Authorization", "Bearer "+accessToken)
	req.URL.RawQuery = query.Encode()
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Error(err)
	}
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Wrong status code returned %v\n", resp.StatusCode)
	}
}

func TestUpdatePasswordWrongKey(t *testing.T) {
	defer dropData(t)
	accessToken, err := createAndAuthenticateMaster()
	if err != nil {
		t.Fatal(err)
	}

	// create a new password
	input := `{"key": "somekey", "password":"somepassword"}`
	req, _ := http.NewRequest(http.MethodPost, baseUrl+"/password/store", bytes.NewReader([]byte(input)))
	req.Header.Add("Authorization", "Bearer "+accessToken)
	http.DefaultClient.Do(req)

	// update password
	input = `{"password":"updatepassword"}`
	req, _ = http.NewRequest(http.MethodPut, baseUrl+"/password/update", bytes.NewReader([]byte(input)))
	query := req.URL.Query()
	query.Add("key", "wrongkey")
	req.Header.Add("Authorization", "Bearer "+accessToken)
	req.URL.RawQuery = query.Encode()
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Error(err)
	}
	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("Wrong status code returned %v\n", resp.StatusCode)
	}
}

func TestRetrieveKeys(t *testing.T) {
	defer dropData(t)
	accessToken, err := createAndAuthenticateMaster()
	if err != nil {
		t.Fatal(err)
	}

	// create a new password
	input := `{"key": "somekey", "password":"somepassword"}`
	req, _ := http.NewRequest(http.MethodPost, baseUrl+"/password/store", bytes.NewReader([]byte(input)))
	req.Header.Add("Authorization", "Bearer "+accessToken)
	http.DefaultClient.Do(req)

	req, err = http.NewRequest(http.MethodGet, baseUrl+"/password/keys", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("Authorization", "Bearer "+accessToken)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	keys, err := readResponse[[]string](resp)
	if err != nil {
		t.Fatal(err)
	}

	if len(*keys) != 1 {
		t.Fatalf("Wrong number of keys returned %v\n", len(*keys))
	}

	if (*keys)[0] != "somekey" {
		t.Fatalf("Wrong key returned %v\n", (*keys)[0])
	}
}
