package iputils

import "testing"

func TestUtils(t *testing.T) {
	// Check if string is valid ip
	if !NameIsIP("127.0.0.1") {
		t.Fatal("Ip address 127.0.0.1 is rendered invalid")
	}
	if NameIsIP("test") {
		t.Fatal("Non ip address is rendered valid")
	}
}
