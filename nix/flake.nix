{
  description = "gofishies flake";
  inputs.flake-utils.url = "github:numtide/flake-utils";
  outputs =
    {
      self,
      nixpkgs,
      flake-utils,
    }:
    flake-utils.lib.eachDefaultSystem (
      system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
      in
      {
        packages = {
          default =
            with pkgs;
            buildGoModule {
              pname = "gofishies";
              version = "0.0.1";

              src = ../.;
              vendorHash = "sha256-8UprJXRLFO3giWAm8k+vbNz7HPYwKW7cD36qc3hEkzE=";

              env.CGO_ENABLED = 0;
            };

          # wasm =
          #   with pkgs;
          #   stdenv.mkDerivation {
          #     name = "gofishies-wasm";
          #     nativeBuildInputs = [ go ];
          #     src = ../.;
          #     buildPhase = ''
          #       GOOS=js GOARCH=wasm go build -o $out/bin/gofishies.wasm
          #     '';
          #   };
        };

        devShells.default = pkgs.mkShell {
          packages = with pkgs; [
            go
            gopls
            golangci-lint
            (pkgs.writeShellScriptBin "build-wasm" ''
              GOOS=js GOARCH=wasm go build -o gofishies.wasm
            '')
          ];
        };
      }
    );
}
