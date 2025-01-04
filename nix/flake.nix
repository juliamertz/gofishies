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
        packages.default =
          with pkgs;
          buildGoModule {
            pname = "gofishies";
            version = "0.0.1";

            src = ../.;
            vendorHash = "sha256-8UprJXRLFO3giWAm8k+vbNz7HPYwKW7cD36qc3hEkzE=";

            meta = with lib; {
              description = "";
              homepage = "";
              license = licenses.mit;
            };

            env.CGO_ENABLED = 0;
          };
        devShells.default = pkgs.mkShell {
          packages = with pkgs; [
            go
            gopls
            golangci-lint
          ];
        };
      }
    );
}
