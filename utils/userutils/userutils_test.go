package userutils

import "testing"

func TestUserUtils(t *testing.T) {
	// Test WhosGroupSuperior
	if WhosGroupSuperior("*", "moderator") != "moderator" {
		t.Fatalf("Something's wrong in function WhosGroupSuperior")
	} else {
		t.Log("WhosGroupSuperior: OK!")
	}
}
