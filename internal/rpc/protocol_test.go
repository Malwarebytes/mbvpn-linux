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
