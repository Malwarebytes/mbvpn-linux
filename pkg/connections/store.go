package connections

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

	"github.com/malwarebytes/mbvpn-linux/pkg/config"
)

type Connection struct {
	InterfaceName   string `json:"interface_name"`
	ServerName      string `json:"server_name"`
	DevicePublicKey string `json:"device_public_key"`
	PeerPublicKey   string `json:"peer_public_key"`
}

type Store interface {
	List() ([]Connection, error)
	Save(Connection) error
	Delete(string) error
}

type fileStore struct{ directory config.DirectoryProvider }

func NewStore(directory config.DirectoryProvider) Store { return fileStore{directory: directory} }

func (s fileStore) List() ([]Connection, error) {
	path, err := s.path()
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var connections []Connection
	return connections, json.Unmarshal(data, &connections)
}

func (s fileStore) Save(connection Connection) error {
	connections, err := s.List()
	if err != nil {
		return err
	}
	for index := range connections {
		if connections[index].InterfaceName == connection.InterfaceName {
			connections[index] = connection
			return s.write(connections)
		}
	}
	return s.write(append(connections, connection))
}

func (s fileStore) Delete(name string) error {
	connections, err := s.List()
	if err != nil {
		return err
	}
	remaining := connections[:0]
	for _, connection := range connections {
		if connection.InterfaceName != name {
			remaining = append(remaining, connection)
		}
	}
	return s.write(remaining)
}

func (s fileStore) path() (string, error) {
	dir, err := s.directory.GetConfigDir()
	return filepath.Join(dir, "connections.json"), err
}

func (s fileStore) write(connections []Connection) error {
	path, err := s.path()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return err
	}
	data, err := json.Marshal(connections)
	if err != nil {
		return err
	}
	file, err := os.CreateTemp(filepath.Dir(path), ".connections-*")
	if err != nil {
		return err
	}
	temporary := file.Name()
	defer os.Remove(temporary)
	if err := file.Chmod(0o600); err != nil {
		file.Close()
		return err
	}
	if _, err := file.Write(data); err != nil {
		file.Close()
		return err
	}
	if err := file.Close(); err != nil {
		return err
	}
	return os.Rename(temporary, path)
}
