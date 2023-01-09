{ nixpkgs ? import <nixpkgs> { } }:
let
  pkgs = {
    atomi = (
      with import (fetchTarball "https://github.com/kirinnee/test-nix-repo/archive/refs/tags/v15.3.0.tar.gz");
      {
        inherit pls spacectl upstash;
      }
    );
    "Unstable 4th July 2022" = (
      with import (fetchTarball "https://github.com/NixOS/nixpkgs/archive/0ea7a8f1b939d74e5df8af9a8f7342097cdf69eb.tar.gz") { };
      {
        inherit
          coreutils
          gnugrep
          bash
          jq

          pre-commit
          nixpkgs-fmt
          shfmt
          shellcheck

          terraform
          terraform-docs
          infracost
          tfsec
          tflint

          kubectx
          ;
        prettier = nodePackages.prettier;
      }
    );

    "Unstable 15th August 2022" = (
      with import (fetchTarball "https://github.com/NixOS/nixpkgs/archive/bc4b4a50c7a105c56f1b712a87818678298deef3.tar.gz") { };
      {
        inherit
          awscli2
          kubectl
          kustomize
          kubelogin-oidc
          findutils
          git;
      }
    );

    "Unstable 14th September 2022" = (
      with import (fetchTarball "https://github.com/NixOS/nixpkgs/archive/ee01de29d2f58d56b1be4ae24c24bd91c5380cea.tar.gz") { };
      {
        inherit
          _1password
          go;
      }
    );

  };
in

with pkgs;

pkgs.atomi // pkgs."Unstable 4th July 2022" // pkgs."Unstable 15th August 2022" // pkgs."Unstable 14th September 2022"
