#!/bin/bash

ROOT_PATH=$(git -C "$(dirname "$0")" rev-parse --show-toplevel)
export ROOT_PATH

function running_in_windows {
  [ "$OS" = "Windows_NT" ]
}
function running_in_gitbash {
  [ "$MSYSTEM" = "MINGW64" ] || [ "$MSYSTEM" = "MINGW32" ] || [ "$MSYSTEM" = "MSYS" ]
}
function is_go_installed {
  command -v go &> /dev/null
}
function is_docker_installed {
  command -v docker &> /dev/null
}

# get version from go.mod file
GO_MOD_VERSION="$(grep -oP 'go \K(\d+\.\d+)' < "$ROOT_PATH"/go.mod)"
export GO_MOD_VERSION

# get version from go command
if is_go_installed; then
    GO_INSTALLED_VERSION="$(go version | grep -oP 'go version go\K(\d+\.\d+)')"
    export GO_INSTALLED_VERSION
fi
