let
  nixpkgs = fetchTarball "https://github.com/NixOS/nixpkgs/tarball/nixos-unstable";
  pkgs = import nixpkgs { config = { }; overlays = [ ]; };
  build = pkgs.callPackage ./build.nix { };
  HOME = builtins.getEnv "HOME";
  PROJECT_ROOT = builtins.toString ./.;
in
{
  build = build;
  shell = pkgs.mkShell {
    hardeningDisable = [ "fortify" ];
    inputsFrom = [ build ];
    packages = [

      pkgs.git

      pkgs.nixpkgs-fmt

    ];

    shellHook = ''
      goversion=$(go version)
      export GOPATH=$HOME/gohome/"$goversion"/go
      export GOCACHE=$HOME/gohome/"$goversion"/cache
      export GOENV=$HOME/gohome/"$goversion"/env
      export PATH=$GOPATH/bin:$PATH
    '';

  };
}
