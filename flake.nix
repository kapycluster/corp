{
  inputs.nixpkgs.url = "github:nixos/nixpkgs";
  inputs.templ.url = "github:a-h/templ";

  outputs = { self, nixpkgs, templ, ... }:
    let
      supportedSystems = [ "x86_64-linux" "aarch64-linux" ];
      forAllSystems = nixpkgs.lib.genAttrs supportedSystems;
      nixpkgsFor = forAllSystems (system: import nixpkgs { inherit system; });

      builder = { pkgs, pname, src, subPackages, enableStatic ? false }: pkgs.buildGoModule {
        inherit pname src;
        version = "1.0.0";

        # Replace with pkgs.lib.fakeHash if you bump go.mod and paste the
        # resulting 'got' back here.
        vendorHash = "sha256-zTqAqojXGrf3gAhDsFxZOKSV9WRpk64fA91LIGxsdm8=";
        proxyVendor = true;
        doCheck = false;
        subPackages = subPackages;

        # Disable CGO for all builds
        CGO_ENABLED = "0";

        # Conditional build inputs and flags based on static linking requirement
        nativeBuildInputs = pkgs.lib.optionals enableStatic [ pkgs.musl ];
        ldflags = pkgs.lib.optionals enableStatic [ "-s" "-w" ''-extldflags "-static -L${pkgs.musl}/lib"'' ];

        postInstall = ''
          mv $out/bin/cmd $out/bin/${pname}
        '';
      };


    in
    {
      packages = forAllSystems (system:
        let
          pkgs = nixpkgsFor.${system};
        in
        {
          controller = builder {
            pkgs = pkgs;
            pname = "controller";
            src = pkgs.lib.cleanSource ./.;
            subPackages = [ "controller/cmd" ];
          };

          kapyserver = builder {
            pkgs = pkgs;
            pname = "kapyserver";
            src = pkgs.lib.cleanSource ./.;
            subPackages = [ "kapyserver/cmd" ];
            enableStatic = true; # Enable static linking for kapyserver
          };

          panel = builder {
            pkgs = pkgs;
            pname = "panel";
            src = pkgs.lib.cleanSource ./.;
            subPackages = [ "panel/cmd" ];
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
              cni-plugins
              templ.packages.${system}.templ
            ];
          };

          apps = forAllSystems (
            system:
            let
              pkgs = nixpkgsFor.${system};
            in
            {
              panel = {
                type = "app";
                program = "${self.packages.${system}.panel}/bin/panel";
                cwd = ./.;
              };
            }
          );
        });
    };
}
