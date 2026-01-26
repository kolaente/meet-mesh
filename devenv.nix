{ pkgs, lib, config, inputs, ... }:

let
  unstable = import inputs.nixpkgs-unstable { system = pkgs.stdenv.system; };
  beads = pkgs.buildGoModule {
    pname = "beads";
    version = "0.47.1";

    src = pkgs.fetchFromGitHub {
      owner = "steveyegge";
      repo = "beads";
      rev = "v0.47.1";
      hash = "sha256-DwIR/r1TJnpVd/CT1E2OTkAjU7k9/KHbcVwg5zziFVg=";
    };

    vendorHash = "sha256-pY5m5ODRgqghyELRwwxOr+xlW41gtJWLXaW53GlLaFw=";

    subPackages = [ "cmd/bd" ];

    ldflags = [ "-s" "-w" ];

    # Tests require git which isn't available in the Nix sandbox
    doCheck = false;
  };
in
{
  packages = with pkgs; [
    beads
  ];

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
