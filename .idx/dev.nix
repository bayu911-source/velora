{ pkgs, ... }: {
  # Let Go applications use the network.
  previews = [
    {
      command = ["go", "run", "velora/cmd/server/main.go"];
      manager = "web";
      port = 8080;
    }
  ];

  # Needed for running Go.
  packages = [
    pkgs.go
  ];
}
