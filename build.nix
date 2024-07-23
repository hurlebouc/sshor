{ lib, buildGoModule, nix-gitignore }:

buildGoModule {
  pname = "sshor-go";
  version = "1.0.1";
  #doCheck = false;

  src = nix-gitignore.gitignoreSource [ ".git" ".gitignore" "*.nix" ] ./.;

  vendorHash = "sha256-CHOhMeT3fN/eebTP0DuzwawMPBlnekYa0ebCODYWh5A="; # get hash after first build
}
