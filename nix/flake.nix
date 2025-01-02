{
  description = "Flake utils demo";

  inputs.flake-utils.url = "github:numtide/flake-utils";

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let pkgs = nixpkgs.legacyPackages.${system}; in
      {
        # packages = rec {
        #   hello = pkgs.hello;
        #   default = hello;
        # };
        devShells.default = pkgs.mkShell {
          packages = with pkgs; [ go gopls ];
        };
      }
    );
}
