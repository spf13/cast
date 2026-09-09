{ pkgs, ... }:

{
  languages = {
    go = {
      enable = true;
      package = pkgs.go_1_27;
    };
  };

  packages = with pkgs; [
    just
    golangci-lint
  ];
}
