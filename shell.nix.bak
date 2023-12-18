{ pkgs ? import <nixpkgs> {} }:
  pkgs.mkShell {
    # nativeBuildInputs is usually what you want -- tools you need to run
    nativeBuildInputs = with pkgs.buildPackages; [
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
      nodejs_21
      sqlite
      rsync
      nomad
      vault
      dig
      vegeta
      just
    ];

    LOCALE_ARCHIVE = "${pkgs.glibcLocales}/lib/locale/locale-archive";
}