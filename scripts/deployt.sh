#!/bin/bash


make clean
make build 


zip -r goplay.zip ./bin 
zip -r goplay.zip ./etc
zip -j goplay.zip ./scripts/install.sh
