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
          controller = self.packages.${system}.controller;
        in
        {
          panel = pkgs.buildGo122Module {
            name = "panel";
            rev = "master";
            src = ./panel/cmd;

            vendorHash = "sha256-EBVD/RzVpxNcwyVHP1c4aKpgNm4zjCz/99LvfA0Oc/Q=";
          };
          panelContainer = pkgs.dockerTools.buildLayeredImage {
            name = "ghcr.io/kapycluster/panel";
            tag = "latest";
            config = {
              Cmd = "${panel}/bin/panel";
              ExposedPorts = { "8080/tcp" = { }; };
            };
          };

          controller = pkgs.buildGo122Module {
            name = "controller";
            rev = "master";
            src = ./controller/cmd;
            vendorHash = "sha256-EBVD/RzVpxNcwyVHP1c4aKpgNm4zjCz/99LvfA0Oc/Q=";
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
              kubectl
              kubectx
              protobuf
              protoc-gen-go
              protoc-gen-go-grpc
            ];
          };
        });
    };
}
