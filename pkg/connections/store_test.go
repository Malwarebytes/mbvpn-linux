package connections

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/malwarebytes/mbvpn-linux/pkg/config"
)

func TestStoreLifecycle(t *testing.T) {
	store := NewStore(config.NewDirectoryProvider(t.TempDir()))
	connection := Connection{InterfaceName: "mbv123", ServerName: "server", DevicePublicKey: "device", PeerPublicKey: "peer"}
	if err := store.Save(connection); err != nil {
		t.Fatalf("Save() error = %v", err)
	}
	connections, err := store.List()
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}
	if len(connections) != 1 || connections[0] != connection {
		t.Fatalf("List() = %#v", connections)
	}
	if err := store.Delete(connection.InterfaceName); err != nil {
		t.Fatalf("Delete() error = %v", err)
	}
	connections, err = store.List()
	if err != nil || len(connections) != 0 {
		t.Fatalf("List() after Delete = %#v, %v", connections, err)
	}
}

func TestStoreUsesPrivateFileMode(t *testing.T) {
	baseDir := t.TempDir()
	store := NewStore(config.NewDirectoryProvider(baseDir))
	if err := store.Save(Connection{InterfaceName: "mbv123"}); err != nil {
		t.Fatalf("Save() error = %v", err)
	}
	info, err := os.Stat(filepath.Join(baseDir, "mbvpn", "connections.json"))
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm() != 0o600 {
		t.Fatalf("file mode = %o, want 600", info.Mode().Perm())
	}
}
