{ pkgs, ... }: {
  # Needed for running Go and related tools.
  packages = [
    pkgs.go
    pkgs.gopls
    pkgs.golangci-lint
  ];

  # Add Go tools to the PATH
  env.PATH = pkgs.lib.makeBinPath [
    pkgs.go
    pkgs.gopls
    pkgs.golangci-lint
  ];
}
