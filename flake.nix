{
  description = "build squishy";

  inputs = {
    nixpkgs.url = "github:kompismoln/nixpkgs/nixos-unstable";
  };

  outputs =
    {
      self,
      nixpkgs,
    }:
    let
      pname = "squishy";
      version = builtins.readFile ./.version;
      src = self;
      systems = [
        "x86_64-linux"
        "aarch64-linux"
      ];
      forAllSystems = f: nixpkgs.lib.genAttrs systems (system: f system);
    in
    {
      packages = forAllSystems (
        system:
        let
          pkgs = nixpkgs.legacyPackages.${system};
        in
        {
          default = pkgs.buildGoModule {
            inherit src pname version;
            vendorHash = "sha256-whS5168qFVPAt+PqeO010YFGR4VG9nABlKqz/eJm/Sk=";
            ldflags = [ "-X main.Version=${version}" ];
          };
        }
      );

      devShells = forAllSystems (
        system:
        let
          pkgs = nixpkgs.legacyPackages.${system};
        in
        {
          default = pkgs.mkShell {
            name = "${pname}-dev";
            packages = with pkgs; [
              air
              go
            ];
          };
        }
      );
    };
}
