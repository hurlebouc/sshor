{ lib, buildGoModule, nix-gitignore }:

buildGoModule {
  pname = "sshor-go";
  version = "1.0.4";
  #doCheck = false;

  src = nix-gitignore.gitignoreSource [ ".git" ".gitignore" "*.nix" ] ./.;

  vendorHash = "sha256-BK4bn/XSIDskmz4ksTmYforGPcOzkp2RPr3sgAGKjQ8="; # get hash after first build
}
