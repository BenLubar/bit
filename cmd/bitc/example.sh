#!/bin/sh

set -e

test -f hello.bit.xz || wget https://benlubar.github.io/useless-crap/hello.bit.xz
(cd ../bit2bit && go build)
go build
xzcat hello.bit.xz | ./bitc -o hello
./hello | ../bit2bit/bit2bit -d
