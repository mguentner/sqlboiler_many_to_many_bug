{
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
  };
  outputs = { self, nixpkgs }:
    let
      supportedSystems = [ "x86_64-linux" "aarch64-linux" ];

      forAllSystems = nixpkgs.lib.genAttrs supportedSystems;

      nixpkgsFor = forAllSystems (system: import nixpkgs {
        inherit system;
        overlays = [ ];
      });
    in
    {
      devShells = forAllSystems (system:
        let
          pkgs = nixpkgsFor.${system};
        in
        {
          default = pkgs.mkShell {
            name = "default";
            buildInputs = with pkgs; [
              go
              gopls
              bashInteractive
              sqlboiler
              postgresql
            ];
            env = {
              "PGHOST" = "127.0.0.1";
              "PGPORT" = "5432";
              "PGUSER" = "postgres";
              "PGPASSWORD" = "secret";
              "PGDATABASE" = "sqlboiler_bug";
            };
          };
        });

    };
}
