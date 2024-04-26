package services

import (
	"testing"

	"github.com/danilomarques1/fidusserver/dtos"
)

func TestRegisterService(t *testing.T) {
	cases := []struct {
		label    string
		input    *dtos.CreateMasterDto
		expected error
	}{

		{"Should register a new master", &dtos.CreateMasterDto{Name: "Mocked Name", Email: "mock@email.com", Password: "thisisasecretpassword"}, nil},
	}

	for _, tc := range cases {
		t.Run(tc.label, t testing.T) {
		}
	}
}
