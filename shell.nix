let nixpkgs = import <nixpkgs> { }; in
let shells = import ./nix/shells.nix { inherit nixpkgs; }; in
shells.dev
