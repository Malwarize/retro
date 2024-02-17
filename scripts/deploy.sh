#!/bin/bash

make clean
make build 

./scripts/compile_installer.sh
zip ./bin/goplay_installer.zip ./bin/install.sh -j
