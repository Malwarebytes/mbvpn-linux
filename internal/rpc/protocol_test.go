package rpc

import (
	"bytes"
	"encoding/binary"
	"testing"
)

func TestRequestRoundTrip(t *testing.T) {
	var buffer bytes.Buffer
	request := Request{ID: "one", Method: "status"}
	if err := WriteRequest(&buffer, request); err != nil {
		t.Fatalf("WriteRequest() error = %v", err)
	}
	got, err := ReadRequest(&buffer)
	if err != nil {
		t.Fatalf("ReadRequest() error = %v", err)
	}
	if got.ID != request.ID || got.Method != request.Method {
		t.Fatalf("ReadRequest() = %#v", got)
	}
}

func TestReadRequestRejectsOversizedFrame(t *testing.T) {
	var buffer bytes.Buffer
	if err := binary.Write(&buffer, binary.BigEndian, uint32(MaxFrameSize+1)); err != nil {
		t.Fatal(err)
	}
	if _, err := ReadRequest(&buffer); err == nil {
		t.Fatal("expected oversized frame error")
	}
}

func TestErrorMessage(t *testing.T) {
	err := (&Error{Code: "not_authenticated", Message: "there is no active session"}).Error()
	if err != "not_authenticated: there is no active session" {
		t.Fatalf("unexpected error message: %q", err)
	}
}

func TestReadRequestRejectsUnknownFields(t *testing.T) {
	var buffer bytes.Buffer
	payload := []byte(`{"version":1,"id":"one","method":"status","extra":true}`)
	if err := binary.Write(&buffer, binary.BigEndian, uint32(len(payload))); err != nil {
		t.Fatal(err)
	}
	if _, err := buffer.Write(payload); err != nil {
		t.Fatal(err)
	}
	if _, err := ReadRequest(&buffer); err == nil {
		t.Fatal("expected unknown field error")
	}
}
