package entry

import (
	"bytes"
	"testing"
)

func TestNewLog(t *testing.T) {
	index := int64(1)
	sequence := 10
	payload := []byte("test payload")

	log := NewLog(index, sequence, payload)

	if log.Index != index {
		t.Errorf("NewLog().Index = %v, want %v", log.Index, index)
	}
	if log.Sequence != sequence {
		t.Errorf("NewLog().Sequence = %v, want %v", log.Sequence, sequence)
	}
	if !bytes.Equal(log.PayLoad, payload) {
		t.Errorf("NewLog().PayLoad = %v, want %v", log.PayLoad, payload)
	}
}

func TestEncodeLog(t *testing.T) {
	log := Log{
		Index:    1,
		Sequence: 10,
		PayLoad:  []byte("test payload"),
	}

	expected := []byte("test payload")

	result := EncodeLog(log)
	if !bytes.Equal(result, expected) {
		t.Errorf("EncodeLog() = %v, want %v", result, expected)
	}
}

func TestDecodeLog(t *testing.T) {
	data := []byte("test payload")

	expected := Log{
		PayLoad: data,
	}

	result, err := DecodeLog(data)
	if err != nil {
		t.Fatalf("DecodeLog() error = %v", err)
	}
	if !bytes.Equal(result.PayLoad, expected.PayLoad) {
		t.Errorf("DecodeLog() = %v, want %v", result.PayLoad, expected.PayLoad)
	}
}
