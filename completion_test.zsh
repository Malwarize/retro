#!/bin/bash
go build -o bin/goplay client/main.go
go build -o bin/goplayer server/main.go
s=`bin/goplay completion zsh`
echo $s > _goplay
source _goplay
