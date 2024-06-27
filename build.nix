{ lib, buildGoModule, nix-gitignore }:

buildGoModule {
  pname = "sshor-go";
  version = "0.1.0";

  src = nix-gitignore.gitignoreSource [".git" ".gitignore" "*.nix"] ./.;

  vendorHash = null; # get hash after first build
}