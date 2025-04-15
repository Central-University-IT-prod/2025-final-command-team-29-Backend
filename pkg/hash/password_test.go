package hash

import (
	"testing"
)

func TestGenerateHash(t *testing.T) {
	password := "mysecretpassword"
	hash, err := GenerateHash(password)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if hash == "" {
		t.Fatalf("Expected a hash, got an empty string")
	}
}

func TestComparePassword(t *testing.T) {
	password := "mysecretpassword"
	hash, err := GenerateHash(password)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	err = ComparePassword(password, hash)
	if err != nil {
		t.Fatalf("Expected passwords to match, got error %v", err)
	}

	wrongPassword := "wrongpassword"
	err = ComparePassword(wrongPassword, hash)
	if err == nil {
		t.Fatalf("Expected passwords to not match, but they did")
	}
}
