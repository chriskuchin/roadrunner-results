{ pkgs ? import <nixpkgs> {} }:
  pkgs.mkShell {
    # nativeBuildInputs is usually what you want -- tools you need to run
    nativeBuildInputs = with pkgs.buildPackages; [
      flyctl
      git
      go
      gotools
      gopls
      go-outline
      gocode
      gopkgs
      gocode-gomod
      godef
      golint
      gcc
      nodejs_18
      sqlite
      rsync
      nomad
    ];

    NOMAD_ADDR=https://n.home.cksuperman.com;
}