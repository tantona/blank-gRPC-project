#!/bin/sh

echo "Compiling proto file"
$GOPATH/src/PROJECT_ROOT/scripts/pb-compile.sh

while inotifywait -r -e modify $GOPATH/src/PROJECT_ROOT/proto; do
    echo “PROTO FILE CHANGED RECOMPILING”
    $GOPATH/src/PROJECT_ROOT/scripts/pb-compile.sh
done
