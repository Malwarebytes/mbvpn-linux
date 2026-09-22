package daemon

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Connection struct {
	ID              string `json:"id"`
	OwnerUID        uint32 `json:"owner_uid"`
	InterfaceName   string `json:"interface_name"`
	ServerName      string `json:"server_name"`
	DevicePublicKey string `json:"device_public_key"`
	PeerPublicKey   string `json:"peer_public_key"`
}

type ConnectionStore struct {
	Paths Paths
}

func (s ConnectionStore) List(uid uint32) ([]Connection, error) {
	path := s.path(uid)
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var connections []Connection
	if err := json.Unmarshal(data, &connections); err != nil {
		return nil, err
	}
	return connections, nil
}

func (s ConnectionStore) Save(uid uint32, connection Connection) error {
	connections, err := s.List(uid)
	if err != nil {
		return err
	}
	updated := false
	for i := range connections {
		if connections[i].ID == connection.ID {
			connections[i] = connection
			updated = true
		}
	}
	if !updated {
		connections = append(connections, connection)
	}
	return s.write(uid, connections)
}

func (s ConnectionStore) Delete(uid uint32, id string) error {
	connections, err := s.List(uid)
	if err != nil {
		return err
	}
	remaining := connections[:0]
	for _, connection := range connections {
		if connection.ID != id {
			remaining = append(remaining, connection)
		}
	}
	return s.write(uid, remaining)
}

func (s ConnectionStore) path(uid uint32) string {
	return filepath.Join(s.Paths.UserConfigDir(uid), "connections.json")
}

func (s ConnectionStore) write(uid uint32, connections []Connection) error {
	dir := s.Paths.UserConfigDir(uid)
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return err
	}
	if err := os.Chmod(dir, 0o700); err != nil {
		return err
	}
	data, err := json.Marshal(connections)
	if err != nil {
		return err
	}
	temporary, err := os.CreateTemp(dir, ".connections-*")
	if err != nil {
		return err
	}
	temporaryPath := temporary.Name()
	defer os.Remove(temporaryPath)
	if err := temporary.Chmod(0o600); err != nil {
		temporary.Close()
		return err
	}
	if _, err := temporary.Write(data); err != nil {
		temporary.Close()
		return err
	}
	if err := temporary.Close(); err != nil {
		return err
	}
	return os.Rename(temporaryPath, s.path(uid))
}

func (s ConnectionStore) Active() (Connection, bool, error) {
	users, err := os.ReadDir(filepath.Join(s.Paths.StateDir, "users"))
	if os.IsNotExist(err) {
		return Connection{}, false, nil
	}
	if err != nil {
		return Connection{}, false, err
	}
	for _, user := range users {
		var uid uint32
		if _, err := fmt.Sscan(user.Name(), &uid); err != nil {
			continue
		}
		connections, err := s.List(uid)
		if err != nil {
			return Connection{}, false, err
		}
		if len(connections) > 0 {
			return connections[0], true, nil
		}
	}
	return Connection{}, false, nil
}
