package servers

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"

	"github.com/Malwarebytes/mbvpn/pkg/remote"
)

const (
	serversFile = "servers.json"
)

type ServerStorage interface {
	Save(locations *remote.VpnLocations) error
	Get() (*remote.VpnLocations, error)
	GetByServerName(name string) (*remote.Server, error)
}

func NewDefaultServerStorage() ServerStorage {
	return &DefaultServerStorage{}
}

type DefaultServerStorage struct{}

func (s *DefaultServerStorage) Save(locations *remote.VpnLocations) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	configDir := filepath.Join(home, ".config", "mbvpn")

	path := filepath.Join(configDir, serversFile)

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	return encoder.Encode(locations)
}

func (s *DefaultServerStorage) Get() (*remote.VpnLocations, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Panic(err)
	}

	configDir := filepath.Join(home, ".config", "mbvpn")

	path := filepath.Join(configDir, serversFile)

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var locations remote.VpnLocations
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&locations)
	if err != nil {
		return nil, err
	}

	return &locations, nil
}

func (s *DefaultServerStorage) GetByServerName(name string) (*remote.Server, error) {
	locations, err := s.Get()
	if err != nil {
		return nil, err
	}

	for _, country := range locations.Countries {
		for _, city := range country.Cities {
			for _, server := range city.Servers {
				if server.Hostname == name {
					return &server, nil
				}
			}
		}
	}

	return nil, nil
}
