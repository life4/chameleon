#!/usr/bin/env bash

echo "make binary releases"
python3 build.py .

echo "make archive"
tar --create --file chameleon.tar \
    --transform='flags=r;s|config_example|config|' \
    --transform='flags=r;s|builds/||' \
    builds templates config_example.toml chameleon.service

echo "delete empty directory from archive"
tar --delete --file chameleon.tar builds

echo "compress"
gzip chameleon.tar

echo "remove binary releases"
rm -r builds
