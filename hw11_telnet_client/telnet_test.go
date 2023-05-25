package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestTelnetClient(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		l, err := net.Listen("tcp", "127.0.0.1:")
		require.NoError(t, err)
		defer func() { require.NoError(t, l.Close()) }()

		var wg sync.WaitGroup
		wg.Add(2)

		go func() {
			defer wg.Done()

			in := &bytes.Buffer{}
			out := &bytes.Buffer{}

			timeout, err := time.ParseDuration("10s")
			require.NoError(t, err)

			client := NewTelnetClient(l.Addr().String(), timeout, io.NopCloser(in), out)
			require.NoError(t, client.Connect())
			defer func() { require.NoError(t, client.Close()) }()

			in.WriteString("hello\n")
			err = client.Send()
			require.NoError(t, err)

			err = client.Receive()
			require.NoError(t, err)
			require.Equal(t, "world\n", out.String())
		}()

		go func() {
			defer wg.Done()

			conn, err := l.Accept()
			require.NoError(t, err)
			require.NotNil(t, conn)
			defer func() { require.NoError(t, conn.Close()) }()

			request := make([]byte, 1024)
			n, err := conn.Read(request)
			require.NoError(t, err)
			require.Equal(t, "hello\n", string(request)[:n])

			n, err = conn.Write([]byte("world\n"))
			require.NoError(t, err)
			require.NotEqual(t, 0, n)
		}()

		wg.Wait()
	})

	t.Run("read from closed connection", func(t *testing.T) {
		in := &bytes.Buffer{}
		out := &bytes.Buffer{}

		timeout, err := time.ParseDuration("10s")
		require.NoError(t, err)

		client := NewTelnetClient("localhost:8080", timeout, io.NopCloser(in), out)

		err = client.Receive()
		require.Truef(t, errors.Is(err, ErrNotConnectionOpen), "actual error %q, excepted %q", err, ErrNotConnectionOpen)
	})

	t.Run("write to closed connection", func(t *testing.T) {
		in := &bytes.Buffer{}
		out := &bytes.Buffer{}

		timeout, err := time.ParseDuration("10s")
		require.NoError(t, err)

		client := NewTelnetClient("localhost:8080", timeout, io.NopCloser(in), out)

		err = client.Send()
		require.Truef(t, errors.Is(err, ErrNotConnectionOpen), "actual error %q, excepted %q", err, ErrNotConnectionOpen)
	})

	t.Run("close closed connection", func(t *testing.T) {
		in := &bytes.Buffer{}
		out := &bytes.Buffer{}

		timeout, err := time.ParseDuration("10s")
		require.NoError(t, err)

		client := NewTelnetClient("localhost:8080", timeout, io.NopCloser(in), out)

		err = client.Close()
		require.Truef(t, errors.Is(err, ErrNotConnectionOpen), "actual error %q, excepted %q", err, ErrNotConnectionOpen)
	})

	t.Run("connect to not exist host", func(t *testing.T) {
		in := &bytes.Buffer{}
		out := &bytes.Buffer{}

		timeout, err := time.ParseDuration("10s")
		require.NoError(t, err)

		client := NewTelnetClient("test.local:8080", timeout, io.NopCloser(in), out)

		err = client.Connect()
		fmt.Println(err)
		var dnsError *net.DNSError
		require.Truef(t, errors.As(err, &dnsError), "actual type error %T, excepted %T", err, dnsError)
	})
}
