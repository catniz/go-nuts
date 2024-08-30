#!/bin/bash
set -e

# Generate all: gen_k8s_yaml.sh dev
# Generate target: gen_k8s_yaml.sh dev deployment.yaml

ENV=$1
TARGET=$2

scripts_dir="$(dirname "$0")"
working_dir="$(realpath "$scripts_dir/../deploy/k8s")"

if [ -z "$ENV" ]; then
  echo "Usage: $0 [ENV] [TARGET]"
  exit 1
fi

function generate_yaml {
  file="$1"
  filename="$(basename "$file")"
  echo "Templating $filename"

  go run "$scripts_dir/gotpl.go" \
    -t "$file" \
    -v "$working_dir/values/$ENV-general.yaml" \
    -v "$working_dir/values/$ENV-infra.yaml" \
    > "$working_dir/yaml/${filename%.tpl}"
}

if [ -z "$TARGET" ]; then
  for file in "$working_dir"/templates/*.tpl; do
    generate_yaml "$file"
  done

else
  file="$working_dir/templates/$TARGET.tpl"
  generate_yaml "$file"
fi

echo "Done"
