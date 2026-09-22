package rpc

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"unicode/utf8"
)

const (
	Version          = 1
	MaxFrameSize     = 64 * 1024
	MaxRequestIDSize = 256
)

var ErrFrameTooLarge = fmt.Errorf("frame exceeds maximum size of %d bytes", MaxFrameSize)

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

func (e *Error) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
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
	if request.Version != Version || request.ID == "" || len(request.ID) > MaxRequestIDSize || request.Method == "" {
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
	if !utf8.Valid(data) {
		return fmt.Errorf("invalid UTF-8 in JSON")
	}
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(target); err != nil {
		return fmt.Errorf("invalid JSON: %w", err)
	}
	var extra any
	if err := decoder.Decode(&extra); err != io.EOF {
		return fmt.Errorf("invalid JSON: trailing data")
	}
	return nil
}

func writeJSON(w io.Writer, value any) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	if len(data) > MaxFrameSize {
		return ErrFrameTooLarge
	}
	if err := binary.Write(w, binary.BigEndian, uint32(len(data))); err != nil {
		return err
	}
	_, err = w.Write(data)
	return err
}
