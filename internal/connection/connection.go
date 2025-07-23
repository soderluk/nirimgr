// Package connection contains functionality to communicate with the Niri socket.
//
// We connect to the Niri socket NIRI_SOCKET, and perform actions and requests to it.
// Actions are same as niri msg action <ACTION>, where <ACTION> is the PascalCase of the
// dashed niri action, e.g. move-window-to-floating -> MoveWindowToFloating etc.
// Requests are the simple requests, e.g. niri msg <REQUEST>, where <REQUEST> is one of the
// niri requests to the socket. niri msg outputs -> Outputs.
// This is heavily inspired by https://github.com/probeldev/niri-float-sticky
package connection

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log/slog"
	"net"
	"os"
	"sync"

	"github.com/soderluk/nirimgr/actions"
	"github.com/soderluk/nirimgr/models"
)

// NiriSocket contains the connection to the socket.
type NiriSocket struct {
	conn net.Conn
}

// pool contains the pool of socket connections.
var pool = sync.Pool{
	New: func() any {
		conn, err := net.Dial("unix", os.Getenv("NIRI_SOCKET"))
		if err != nil {
			slog.Error("Failed to connect to NIRI_SOCKET", "error", err.Error())
			panic(err)
		}
		return &NiriSocket{conn: conn}
	},
}

// Socket returns the NiriSocket from the pool.
func Socket() *NiriSocket {
	sock, ok := pool.Get().(*NiriSocket)
	if !ok {
		slog.Error("Could not get socket")
		panic("could not get socket")
	}
	return sock
}

// PutSocket adds the socket to the pool.
func PutSocket(socket *NiriSocket) {
	pool.Put(socket)
}

// Send writes the request to the socket.
func (s *NiriSocket) Send(req string) error {
	_, err := fmt.Fprintf(s.conn, "%s\n", req)
	return err
}

// Recv reads the data from the socket.
func (s *NiriSocket) Recv() <-chan []byte {
	lines := make(chan []byte)

	go func() {
		defer func() { _ = s.conn.Close() }()
		defer close(lines)

		scanner := bufio.NewScanner(s.conn)
		for scanner.Scan() {
			lines <- scanner.Bytes()
		}

		if err := scanner.Err(); err != nil {
			slog.Error("Could not scan response from socket", "error", err.Error())
		}
	}()

	return lines
}

// Close closes the socket connection.
func (s *NiriSocket) Close() {
	_ = s.conn.Close()
}

// PerformAction performs the given action.
//
// The action is one of the actions that niri can handle.
// The supported actions are defined here: https://docs.rs/niri-ipc/latest/niri_ipc/enum.Action.html
func PerformAction(action actions.Action) bool {
	socket := Socket()
	name := action.GetName()
	slog.Debug("PerformAction", "name", name, "action", action)

	// Convert the action to a map.
	actionData, err := structToMap(action)
	if err != nil {
		slog.Error("Could not convert action to map", "error", err.Error())
		return false
	}
	// We need the request as a string to be sent to the socket.
	request, err := structToString(map[string]any{
		"Action": map[string]any{
			name: actionData,
		},
	})
	if err != nil {
		slog.Error("Could not convert action request to string", "error", err.Error())
		return false
	}
	if err := socket.Send(string(request)); err != nil {
		slog.Error("Error sending request", "error", err.Error())
		return false
	}
	return true
}

// PerformRequest sends a simple request to the niri socket.
//
// The request is one of the requests that niri can handle.
// The supported requests are defined here: https://docs.rs/niri-ipc/latest/niri_ipc/enum.Request.html
func PerformRequest(req models.NiriRequest) (<-chan models.Response, error) {
	stream := make(chan models.Response)
	socket := Socket()

	go func() {
		defer PutSocket(socket)
		defer socket.Close()
		defer close(stream)

		for line := range socket.Recv() {
			if len(line) < 2 {
				continue
			}

			var response models.Response

			if err := json.Unmarshal(line, &response); err != nil {
				slog.Error("Error decoding JSON", "error", err.Error())
				continue
			}
			stream <- response
		}
	}()

	if err := socket.Send(fmt.Sprintf("\"%s\"", req)); err != nil {
		return nil, fmt.Errorf("error requesting event stream: %w", err)
	}

	return stream, nil
}

// structToMap converts a go struct to a map.
func structToMap(a any) (map[string]any, error) {
	var m map[string]any
	b, err := json.Marshal(a)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(b, &m)
	return m, err
}

// structToString converts a go struct to a string.
func structToString(a any) (string, error) {
	b, err := json.Marshal(a)
	if err != nil {
		return "", err
	}

	return string(b), nil
}
