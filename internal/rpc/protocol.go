package rpc

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
)

const (
	Version      = 1
	MaxFrameSize = 64 * 1024
)

type Request struct {
	Version int             `json:"version"`
	ID      string          `json:"id"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type Response struct {
	Version int             `json:"version"`
	ID      string          `json:"id"`
	Result  json.RawMessage `json:"result,omitempty"`
	Error   *Error          `json:"error,omitempty"`
}

func ReadRequest(r io.Reader) (Request, error) {
	var request Request
	if err := readJSON(r, &request); err != nil {
		return Request{}, err
	}
	if request.Version != Version || request.ID == "" || request.Method == "" {
		return Request{}, fmt.Errorf("invalid request")
	}
	return request, nil
}

func WriteResponse(w io.Writer, response Response) error {
	response.Version = Version
	return writeJSON(w, response)
}

func ReadResponse(r io.Reader) (Response, error) {
	var response Response
	if err := readJSON(r, &response); err != nil {
		return Response{}, err
	}
	if response.Version != Version || response.ID == "" {
		return Response{}, fmt.Errorf("invalid response")
	}
	return response, nil
}

func WriteRequest(w io.Writer, request Request) error {
	request.Version = Version
	return writeJSON(w, request)
}

func readJSON(r io.Reader, target any) error {
	var length uint32
	if err := binary.Read(r, binary.BigEndian, &length); err != nil {
		return err
	}
	if length == 0 || length > MaxFrameSize {
		return fmt.Errorf("invalid frame size")
	}
	data := make([]byte, length)
	if _, err := io.ReadFull(r, data); err != nil {
		return err
	}
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(target); err != nil {
		return fmt.Errorf("invalid JSON: %w", err)
	}
	if decoder.More() {
		return fmt.Errorf("invalid JSON")
	}
	return nil
}

func writeJSON(w io.Writer, value any) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	if len(data) > MaxFrameSize {
		return fmt.Errorf("frame exceeds maximum size")
	}
	if err := binary.Write(w, binary.BigEndian, uint32(len(data))); err != nil {
		return err
	}
	_, err = w.Write(data)
	return err
}
