{ lib, buildGoModule, nix-gitignore }:

buildGoModule {
  pname = "sshor-go";
  version = "0.9.0";

  src = nix-gitignore.gitignoreSource [ ".git" ".gitignore" "*.nix" ] ./.;

  vendorHash = "sha256-IFRTzrWW9ewMEoYugy5T3h0ioGklZOWCctlF6vI4z54="; # get hash after first build
}
