package util

import (
	"testing"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestGenerateAndParseToken(t *testing.T) {
	secret := "test-secret-for-unit-test"
	id := primitive.NewObjectID()
	token, err := GenerateToken(secret, id, "tester", "student", "active")
	if err != nil {
		t.Fatalf("GenerateToken error: %v", err)
	}
	claims, err := ParseToken(secret, token)
	if err != nil {
		t.Fatalf("ParseToken error: %v", err)
	}
	if claims.UserID != id.Hex() {
		t.Errorf("claims.UserID = %q, want %q", claims.UserID, id.Hex())
	}
	if claims.Role != "student" {
		t.Errorf("claims.Role = %q", claims.Role)
	}
}

func TestParseTokenInvalid(t *testing.T) {
	if _, err := ParseToken("secret", "invalid.token.value"); err == nil {
		t.Error("expected error for invalid token")
	}
}

func TestParseTokenWrongSecret(t *testing.T) {
	secret := "test-secret-for-unit-test"
	id := primitive.NewObjectID()
	token, err := GenerateToken(secret, id, "tester", "admin", "active")
	if err != nil {
		t.Fatalf("GenerateToken error: %v", err)
	}
	if _, err := ParseToken("wrong-secret", token); err == nil {
		t.Error("expected error for wrong secret")
	}
}
