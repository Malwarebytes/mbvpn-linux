package servers

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strings"

	"github.com/malwarebytes/mbvpn-linux/pkg/config"
	"github.com/malwarebytes/mbvpn-linux/pkg/remote"
)

type ServerStorage interface {
	Save(locations *remote.VpnLocations) error
	Get() (*remote.VpnLocations, error)
	GetByServerName(name string) (*remote.Server, error)
}

func NewDefaultServerStorage(dirProvider config.DirectoryProvider) ServerStorage {
	return &DefaultServerStorage{
		dirProvider: dirProvider,
	}
}

type DefaultServerStorage struct {
	dirProvider config.DirectoryProvider
}

func (s *DefaultServerStorage) Save(locations *remote.VpnLocations) error {
	path, err := s.dirProvider.GetServersFile()
	if err != nil {
		return fmt.Errorf("failed to get servers file path: %w", err)
	}

	// Ensure directory exists
	configDir := filepath.Dir(path)
	if err := os.MkdirAll(configDir, 0o755); err != nil {
		return fmt.Errorf("failed to create config directory '%s': %w", configDir, err)
	}

	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create file at path '%s': %w", path, err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	return encoder.Encode(locations)
}

func (s *DefaultServerStorage) Get() (*remote.VpnLocations, error) {
	path, err := s.dirProvider.GetServersFile()
	if err != nil {
		return nil, fmt.Errorf("failed to get servers file path: %w", err)
	}

	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open file at path '%s': %w", path, err)
	}
	defer file.Close()

	var locations remote.VpnLocations
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&locations)
	if err != nil {
		return nil, fmt.Errorf("failed to decode locations from file: %w", err)
	}

	return &locations, nil
}

func (s *DefaultServerStorage) GetByServerName(name string) (*remote.Server, error) {
	locations, err := s.Get()
	if err != nil {
		return nil, fmt.Errorf("failed to get locations: %w", err)
	}

	for _, country := range locations.Countries {
		for _, city := range country.Cities {
			for _, server := range city.Servers {
				shortHost := strings.Split(server.Hostname, ".")[0]
				if server.Hostname == name || shortHost == name {
					return &server, nil
				}
			}
			if city.Code == name || city.Name == name {
				return &city.Servers[randomInt(len(city.Servers))], nil
			}
		}
		if country.Code == name || country.Name == name {
			randomCity := country.Cities[randomInt(len(country.Cities))]
			randomServer := randomCity.Servers[randomInt(len(randomCity.Servers))]
			return &randomServer, nil
		}
	}

	return nil, nil
}

func randomInt(max int) int {
	if max <= 0 {
		return 0
	}
	return rand.Intn(max)
}
