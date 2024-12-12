package index

import (
	"encoding/hex"
	"testing"
)

func TestNewIndex(t *testing.T) {
	index := NewIndex(1, 100, 200)
	if index.Index != 1 {
		t.Errorf("expected Index to be 1, got %d", index.Index)
	}
	if index.MetadataOffset != 100 {
		t.Errorf("expected MetadataOffset to be 100, got %d", index.MetadataOffset)
	}
	if index.MetadataSize != 200 {
		t.Errorf("expected MetadataSize to be 200, got %d", index.MetadataSize)
	}
}

func TestEncodeIndex(t *testing.T) {
	index := NewIndex(1, 100, 200)
	encoded := EncodeIndex(index)
	expected := "01000000000000006400000000000000c8000000"
	if hex.EncodeToString(encoded) != expected {
		t.Errorf("expected %s, got %s", expected, hex.EncodeToString(encoded))
	}
}

func TestDecodeIndex(t *testing.T) {
	data, _ := hex.DecodeString("01000000000000006400000000000000c8000000")
	index, err := DecodeIndex(data)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if index.Index != 1 {
		t.Errorf("expected Index to be 1, got %d", index.Index)
	}
	if index.MetadataOffset != 100 {
		t.Errorf("expected MetadataOffset to be 100, got %d", index.MetadataOffset)
	}
	if index.MetadataSize != 200 {
		t.Errorf("expected MetadataSize to be 200, got %d", index.MetadataSize)
	}
}

func TestDecodeIndex_InvalidSize(t *testing.T) {
	data := []byte{0x01, 0x02}
	_, err := DecodeIndex(data)
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}
