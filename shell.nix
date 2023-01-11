{ pkgs ? import <nixpkgs> { } }:

with pkgs;

mkShell {
    name = "badaas-shell";

    buildInputs = [
    go
    gotools
    gopls
    go-outline
    gocode
    gopkgs
    gocode-gomod
    godef
    golint
    golangci-lint
    ];
}
