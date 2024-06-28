{
  description = "panel: ";

  inputs.nixpkgs.url = "github:nixos/nixpkgs";

  outputs =
    { self
    , nixpkgs
    ,
    }:
    let
      supportedSystems = [ "x86_64-linux" "aarch64-linux" ];
      forAllSystems = nixpkgs.lib.genAttrs supportedSystems;
      nixpkgsFor = forAllSystems (system: import nixpkgs { inherit system; });
    in
    {
      packages = forAllSystems (system:
        let
          pkgs = nixpkgsFor.${system};
          panel = self.packages.${system}.panel;
          files = pkgs.lib.fileset.toSource {
            root = ./.;
            fileset = pkgs.lib.fileset.unions [
            ];
          };
        in
        {
          panel = pkgs.buildGo122Module {
            name = "panel";
            rev = "master";
            src = ./.;

            vendorHash = "sha256-EBVD/RzVpxNcwyVHP1c4aKpgNm4zjCz/99LvfA0Oc/Q=";
          };
          panelContainer = pkgs.dockerTools.buildLayeredImage {
            name = "ghcr.io/decantor/panel";
            tag = "latest";
            contents = [ files ];
            config = {
              Cmd = "${panel}/bin/panel";
              ExposedPorts = { "5555/tcp" = { }; };
            };
          };
        });

      defaultPackage = forAllSystems (system: self.packages.${system}.panel);
      devShells = forAllSystems (system:
        let
          pkgs = nixpkgsFor.${system};
        in
        {
          default = pkgs.mkShell {
            nativeBuildInputs = with pkgs; [
              go
            ];
          };
        });
    };
}
