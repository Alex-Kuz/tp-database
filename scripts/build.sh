#!/bin/bash

main=src/main/main.go
app=server.app

dir=$(pwd)

echo $dir

if [[ ! -d 'vendor' ]]; then
    dep ensure -update
    dep ensure
fi

go build -o ${app} ${main}

md5sum ${app}


# for idx in $(for line in $(cat src/sql/dbscheme.sql | grep -E '[a-z_]+_idx'); do echo $line | grep -E '.*_idx'; done;); do echo "DROP INDEX IF EXISTS $idx;"; done;
