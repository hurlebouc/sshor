{ lib, buildGoModule, nix-gitignore }:

buildGoModule {
  pname = "sshor-go";
  version = "0.1.0";

  src = nix-gitignore.gitignoreSource [ ".git" ".gitignore" "*.nix" ] ./.;

  vendorHash = "sha256-IdjZIoI9IL1IpN5SlK/SUrIN3vLFW/D65ZA+ODBAoqg="; # get hash after first build
}
