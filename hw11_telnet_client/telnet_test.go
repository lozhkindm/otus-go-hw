package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"net"
	"os"
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

			client := NewTelnetClient(l.Addr().String(), timeout, ioutil.NopCloser(in), out)
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
}

func TestMessages(t *testing.T) {
	// http://craigwickesser.com/2015/01/capture-stdout-in-go/
	realIn := os.Stdin
	stdin, fakeIn, err := os.Pipe()
	require.NoError(t, err)
	os.Stdin = fakeIn
	defer func() {
		os.Stdin = realIn
	}()

	realErr := os.Stderr
	stderr, fakeErr, err := os.Pipe()
	require.NoError(t, err)
	os.Stderr = fakeErr
	defer func() {
		os.Stderr = realErr
	}()

	t.Run("connection closed", func(t *testing.T) {
		l, err := net.Listen("tcp", "127.0.0.1:")
		require.NoError(t, err)
		defer func() { require.NoError(t, l.Close()) }()

		var wg sync.WaitGroup
		wg.Add(2)

		go func() {
			defer wg.Done()

			out := &bytes.Buffer{}
			scanner := bufio.NewScanner(stderr)

			timeout, err := time.ParseDuration("10s")
			require.NoError(t, err)

			client := NewTelnetClient(l.Addr().String(), timeout, stdin, out)
			require.NoError(t, client.Connect())

			require.True(t, scanner.Scan())
			require.Equal(t, fmt.Sprintf("...connected to %s", l.Addr().String()), scanner.Text())

			err = client.Receive()
			require.NoError(t, err)

			require.True(t, scanner.Scan())
			require.Equal(t, "...connection was closed by peer", scanner.Text())

		}()

		go func() {
			defer wg.Done()

			conn, err := l.Accept()
			require.NoError(t, err)
			require.NotNil(t, conn)

			require.NoError(t, conn.Close())
		}()

		wg.Wait()
	})

	t.Run("EOF", func(t *testing.T) {
		l, err := net.Listen("tcp", "127.0.0.1:")
		require.NoError(t, err)
		defer func() { require.NoError(t, l.Close()) }()

		var wg sync.WaitGroup
		wg.Add(2)

		go func() {
			defer wg.Done()

			out := &bytes.Buffer{}
			scanner := bufio.NewScanner(stderr)

			timeout, err := time.ParseDuration("10s")
			require.NoError(t, err)

			client := NewTelnetClient(l.Addr().String(), timeout, stdin, out)
			require.NoError(t, client.Connect())

			require.True(t, scanner.Scan())
			require.Equal(t, fmt.Sprintf("...connected to %s", l.Addr().String()), scanner.Text())

			err = client.Send()
			require.NoError(t, err)

			require.True(t, scanner.Scan())
			require.Equal(t, "...EOF", scanner.Text())

		}()

		go func() {
			defer wg.Done()

			conn, err := l.Accept()
			require.NoError(t, err)
			require.NotNil(t, conn)
			defer func() { require.NoError(t, conn.Close()) }()

			require.NoError(t, fakeIn.Close())
		}()

		wg.Wait()
	})
}
