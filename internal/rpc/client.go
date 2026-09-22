package rpc

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"time"
)

type Client struct {
	Socket string
}

func (c Client) Call(ctx context.Context, method string, params any, result any) error {
	data, err := json.Marshal(params)
	if err != nil {
		return err
	}
	dialer := net.Dialer{Timeout: 5 * time.Second}
	conn, err := dialer.DialContext(ctx, "unix", c.Socket)
	if err != nil {
		return fmt.Errorf("connect to mbvpnd: %w", err)
	}
	defer conn.Close()
	if deadline, ok := ctx.Deadline(); ok {
		_ = conn.SetDeadline(deadline)
	} else {
		_ = conn.SetDeadline(time.Now().Add(30 * time.Second))
	}
	request := Request{ID: fmt.Sprintf("%d", time.Now().UnixNano()), Method: method, Params: data}
	if err := WriteRequest(conn, request); err != nil {
		return err
	}
	response, err := ReadResponse(conn)
	if err != nil {
		return err
	}
	if response.ID != request.ID {
		return fmt.Errorf("unexpected response")
	}
	if response.Error != nil {
		return response.Error
	}
	if result == nil || len(response.Result) == 0 {
		return nil
	}
	return json.Unmarshal(response.Result, result)
}
