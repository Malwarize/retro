#!/bin/bash
go build -o bin/retro client/main.go
go build -o bin/retroPlayer server/main.go
s=`bin/retro completion zsh`
echo $s > _retro
source _retro
