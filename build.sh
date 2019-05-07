#!/usr/bin/env bash

set -e
set -o pipefail

SRC_DIR=$(dirname $(realpath $0))
MOD_DIR="${SRC_DIR}"
mkdir -p "${MOD_DIR}/usr/local/share/ca-certificates"
cp -r /usr/local/share/ca-certificates "${MOD_DIR}/usr/local/share/ca-certificates"
docker build -f "$SRC_DIR/Dockerfile" "$MOD_DIR"
