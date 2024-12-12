package crc

import (
	"testing"
)

func TestEncode(t *testing.T) {
	data := []byte("test data")
	expected := uint32(3540561586)

	result := Encode(data)
	if result != expected {
		t.Errorf("Encode() = %v, want %v", result, expected)
	}
}

func TestIsMatch(t *testing.T) {
	data := []byte("test data")
	crc := uint32(3540561586)

	if !IsMatch(data, crc) {
		t.Errorf("IsMatch() = false, want true")
	}

	if IsMatch(data, crc+1) {
		t.Errorf("IsMatch() = true, want false")
	}
}
