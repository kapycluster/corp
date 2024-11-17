{
  inputs.nixpkgs.url = "github:nixos/nixpkgs";
  inputs.templ.url = "github:a-h/templ";

  outputs = { self, nixpkgs, templ, ... }:
    let
      supportedSystems = [ "x86_64-linux" "aarch64-linux" ];
      forAllSystems = nixpkgs.lib.genAttrs supportedSystems;
      nixpkgsFor = forAllSystems (system: import nixpkgs {
        inherit system;
        overlays = [ self.overlays.default ];
      });

      builder = { pkgs, pname, src, subPackages, enableStatic ? false }: pkgs.buildGoModule {
        inherit pname src;
        version = "1.0.0";

        # Replace with pkgs.lib.fakeHash if you bump go.mod and paste the
        # resulting 'got' back here.
        vendorHash =
          let
            hashes = builtins.fromJSON (builtins.readFile ./hashes.json);
          in
          hashes.go;

        proxyVendor = true;
        doCheck = false;
        subPackages = subPackages;

        # Disable CGO for all builds
        CGO_ENABLED = "0";

        # Conditional build inputs and flags based on static linking requirement
        nativeBuildInputs = pkgs.lib.optionals enableStatic [ pkgs.musl ];
        ldflags = pkgs.lib.optionals enableStatic [ "-s" "-w" ''-extldflags "-static -L${pkgs.musl}/lib"'' ];

        preBuild =
          if pname == "panel" then ''
            cp -R ${pkgs.panelNodeModules}/node_modules ./panel/views
            cd ./panel
            ${pkgs.tailwindcss}/bin/tailwindcss -c ./views/tailwind.config.js -i ./views/input.css -o ./views/static/style.css
            cd ..
            ${pkgs.templ}/bin/templ generate
            ls -la ./panel/views
          '' else null;

        postInstall = ''
          mv $out/bin/cmd $out/bin/${pname}
        '';
      };

    in
    {
      overlays.default = final: prev: {
        panelNodeModules = with final;
          stdenv.mkDerivation {
            pname = "panel-node-modules";
            version = "0.0.1";
            # impureEnvVars =
            #   lib.fetchers.proxyImpureEnvVars
            #   ++ [ "GIT_PROXY_COMMAND" "SOCKS_SERVER" ];
            src = ./.;
            nativeBuildInputs = [ bun ];
            buildInputs = [ nodejs-slim_latest ];
            dontConfigure = true;
            dontFixup = true;
            buildPhase = ''
              cd panel/views
              bun i --no-progress --frozen-lockfile
              cd ../..
            '';
            installPhase = ''
              mkdir -p $out
              ls -la
              cp -R ./panel/views/node_modules $out/
            '';

            outputHash =
              let
                hashes = builtins.fromJSON (builtins.readFile ./hashes.json);
              in
              hashes.node;

            outputHashAlgo = "sha256";
            outputHashMode = "recursive";
          };

        controller = builder {
          pkgs = final;
          pname = "controller";
          src = final.lib.cleanSource ./.;
          subPackages = [ "controller/cmd" ];
        };

        kapyserver = builder {
          pkgs = final;
          pname = "kapyserver";
          src = final.lib.cleanSource ./.;
          subPackages = [ "kapyserver/cmd" ];
          enableStatic = true; # Enable static linking for kapyserver
        };

        panel = builder {
          pkgs = final;
          pname = "panel";
          src = final.lib.cleanSource ./.;
          subPackages = [ "panel/cmd" ];
        };
      };

      packages = forAllSystems (system:
        {
          inherit (nixpkgsFor.${system}) controller kapyserver panel panelNodeModules;
        });

      defaultPackage = forAllSystems (system: self.packages.${system}.panel);

      devShells = forAllSystems (system:
        let
          pkgs = nixpkgsFor.${system};
        in
        {
          default = pkgs.mkShell {
            nativeBuildInputs = with pkgs; [
              k3d
              go_1_23
              kubectl
              kubectx
              protobuf
              protoc-gen-go
              protoc-gen-go-grpc
              cni-plugins
              templ.packages.${system}.templ
              tailwindcss
              bun
              nodePackages.nodejs
              gopls
              air
              gcc
              nixd
              kustomize
            ];
          };
        });

      apps = forAllSystems (
        system:
        {
          panel = {
            type = "app";
            program = "${self.packages.${system}.panel}/bin/panel";
            cwd = ./.;
          };
        }
      );

      formatter = forAllSystems (system: nixpkgsFor."${system}".nixpkgs-fmt);
    };
}
