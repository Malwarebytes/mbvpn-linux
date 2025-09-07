{
  description = "MBVPN development environment";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
      in
      {
        devShells.default = pkgs.mkShell {
          buildInputs = with pkgs; [
            go
            wireguard-tools
            git
            gnumake
            podman
            podman-compose
            podman-desktop
            sudo
          ];

          shellHook = ''
            echo "MBVPN development environment"
            echo "Go version: $(go version)"
            echo "WireGuard tools available"
          '';
        };
      });
}
