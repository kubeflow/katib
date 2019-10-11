package env

import (
	"os"
	"testing"
)

func TestGetEnvWithDefault(t *testing.T) {
	expected := "FAKE"
	key := "TEST"
	v := GetEnvOrDefault(key, expected)
	if v != expected {
		t.Errorf("Expected %s, got %s", expected, v)
	}
	expected = "FAKE1"
	os.Setenv(key, expected)
	v = GetEnvOrDefault(key, "")
	if v != expected {
		t.Errorf("Expected %s, got %s", expected, v)
	}
}
