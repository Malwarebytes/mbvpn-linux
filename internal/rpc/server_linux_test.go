//go:build linux

package rpc

import (
	"errors"
	"net"
	"os"
	"path/filepath"
	"sync"
	"testing"
)

func TestPrepareSocket(t *testing.T) {
	t.Run("missing path", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "s")
		err := prepareSocket(path)
		if err != nil {
			t.Fatalf("prepareSocket(): %v", err)
		}
		_, err = os.Stat(path)
		if err == nil || !errors.Is(err, os.ErrNotExist) {
			t.Fatalf("os.Stat(): path still exists, got %v", err)
		}
	})
	t.Run("live listener", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "test.sock")
		ln, err := net.Listen("unix", path)
		if err != nil {
			t.Fatalf("Listen: %v", err)
		}
		defer ln.Close()
		err = prepareSocket(path)
		if !errors.Is(err, ErrSocketInUse) {
			t.Fatalf("prepareSocket(): expected ErrSocketInUse, got %v", err)
		}
		_, err = os.Stat(path)
		if err != nil {
			t.Fatalf("os.Stat(): socket file is missing: %v", err)
		}
	})
	t.Run("stale socket", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "test.sock")
		ln, err := net.Listen("unix", path)
		if err != nil {
			t.Fatalf("Listen: %v", err)
		}
		unixLn := ln.(*net.UnixListener)
		unixLn.SetUnlinkOnClose(false)
		if err := unixLn.Close(); err != nil {
			t.Fatalf("Close: %v", err)
		}
		_, err = os.Stat(path)
		if err != nil {
			t.Fatalf("os.Stat(): socket file is missing: %v", err)
		}
		err = prepareSocket(path)
		if err != nil {
			t.Fatalf("prepareSocket(): %v", err)
		}
		_, err = os.Stat(path)
		if err == nil || !errors.Is(err, os.ErrNotExist) {
			t.Fatalf("os.Stat(): path still exists, got %v", err)
		}
	})
	t.Run("regular file", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "test.t")
		err := os.WriteFile(path, []byte("test"), 0o644)
		if err != nil {
			t.Fatalf("WriteFile: %v", err)
		}
		err = prepareSocket(path)
		if err == nil {
			t.Fatalf("prepareSocket(): expected error, got %v", err)
		}
		if errors.Is(err, ErrSocketInUse) {
			t.Fatalf("prepareSocket(): expected different error, got %v", err)
		}
		_, err = os.Stat(path)
		if err != nil {
			t.Fatalf("os.Stat(): file was removed: %v", err)
		}
	})
	t.Run("two overlapping claims", func(t *testing.T) {
		dir := t.TempDir()
		socketPath := filepath.Join(dir, "s")
		lockPath := filepath.Join(dir, "daemon.lock")

		var mu sync.Mutex
		var listeners []net.Listener
		defer func() {
			mu.Lock()
			defer mu.Unlock()
			for _, ln := range listeners {
				ln.Close()
			}
		}()

		start := make(chan struct{})
		errs := make(chan error, 2)
		claim := func() error {
			return withSocketLock(lockPath, func() error {
				if err := prepareSocket(socketPath); err != nil {
					return err
				}
				ln, err := net.Listen("unix", socketPath)
				if err != nil {
					return err
				}
				mu.Lock()
				listeners = append(listeners, ln)
				mu.Unlock()
				return nil
			})
		}

		for range 2 {
			go func() {
				<-start
				errs <- claim()
			}()
		}
		close(start)

		var succeeded, inUse int
		for range 2 {
			err := <-errs
			if err == nil {
				succeeded++
				continue
			}
			if !errors.Is(err, ErrSocketInUse) {
				t.Fatalf("claim: %v", err)
			}
			inUse++
		}
		if succeeded != 1 || inUse != 1 {
			t.Fatalf("expected 1 success and 1 ErrSocketInUse, got %d success and %d ErrSocketInUse", succeeded, inUse)
		}
		conn, err := net.Dial("unix", socketPath)
		if err != nil {
			t.Fatalf("cannot connect to socket: %v", err)
		}
		conn.Close()
	})
}
