{ pkgs, lib, config, inputs, ... }:

let
  unstable = import inputs.nixpkgs-unstable { system = pkgs.stdenv.system; };
in
{

  languages = {
    javascript = {
      enable = true;
      pnpm.enable = true;
    };

    go = {
      enable = true;
      package = pkgs.go_1_25;
    };
  };

  services = {
    mailhog.enable = true;
  };
}
