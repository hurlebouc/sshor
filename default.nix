let
  nixpkgs = fetchTarball "https://github.com/NixOS/nixpkgs/archive/b73c2221a46c13557b1b3be9c2070cc42cf01eb3.tar.gz";
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
      pkgs.gh
      pkgs.podman-compose
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
