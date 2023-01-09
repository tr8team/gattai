{ nixpkgs ? import <nixpkgs> { } }:
let pkgs = import ./packages.nix { inherit nixpkgs; }; in
with pkgs;
{
  system = [
    findutils
    coreutils
    gnugrep
    jq
    bash
  ];

  dev = [
  ];

  main = [
    awscli2
    kubectl
    kustomize
    kubelogin-oidc
    spacectl
    pls
    git
    terraform
    terraform-docs
    kubectx
    _1password
    go
    upstash
  ];

  lint = [
    pre-commit
    nixpkgs-fmt
    prettier
    shfmt
    shellcheck
    infracost

    tfsec
    tflint
  ];

  ops = [
  ];

}
