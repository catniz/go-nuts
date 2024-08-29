#!/bin/bash

ROOT_PATH=$(git -C "$(dirname "$0")" rev-parse --show-toplevel)
export ROOT_PATH

source "$ROOT_PATH/scripts/make_env.sh"

# fix locale error that may occur when calling grep -P in git bash environment
export LANG=C.UTF-8
export LC_ALL=C.UTF-8

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

function run_in_dir {
  dir=$1
  shift
  pushd "$dir" > /dev/null || (echo "Failed to pushd to $dir" && exit 1)
  "$@"
  popd > /dev/null || (echo "Failed to popd from $dir" && exit 1)
}

function go_cmd {
  if is_go_installed && [ "$GO_MOD_VERSION" = "$GO_INSTALLED_VERSION" ]; then
    # if go version is compatible, build directly
    run_in_dir "$ROOT_PATH" go "$@"

  elif is_docker_installed; then
    if running_in_gitbash; then
      # if running in Windows Git Bash (MINGW), cd to wsl directory to mount volume correctly
      MSYS2_ARG_CONV_EXCL='*' run_in_dir "$ROOT_PATH" docker run --rm -v .:/app:rw -w /app golang:"${GO_MOD_VERSION}"-alpine go "$@"
    else
      # build with docker in normal Windows environment
      docker run --rm -v "$ROOT_PATH":/app:rw -w /app golang:"${GO_MOD_VERSION}"-alpine go "$@"
    fi
  else
    echo "Please install the appropriate version of go or docker, or check the path."
    exit 1
  fi
}