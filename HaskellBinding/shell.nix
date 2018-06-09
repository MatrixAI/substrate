{
  pkgs ? import (fetchTarball https://github.com/NixOS/nixpkgs-channels/archive/7cbf6ca1c84dfc917c1a99524e082fb677501844.tar.gz) {}
}:
  with pkgs;
  haskell.lib.buildStackProject {
    name = "container-demo";
    buildInputs = [go];
    shellHook = ''
      echo 'Building container demo'
      set -v
      alias stack="\stack --nix"
      set +v
      echo 'Setting GOPATH'
      export GOPATH='./hs-libcontainer/src/go'
    '';
  }