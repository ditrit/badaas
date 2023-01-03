{ pkgs ? import <nixpkgs> { } }:

with pkgs;

mkShell {
    name = "badaas-oidc-shell";

    buildInputs = [
        nodejs-18_x
    ];
}
