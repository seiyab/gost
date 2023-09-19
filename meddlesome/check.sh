#!/usr/bin/env -S bash -eu

SCRIPT_DIR=$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" &> /dev/null && pwd)

GIT_ROOT_DIR="$(git rev-parse --show-toplevel)"
cd -- "$GIT_ROOT_DIR/gost"
go install ./

cd -- "$SCRIPT_DIR"

for MODULE in $(ls "$SCRIPT_DIR/objectives")
do
  cd -- "$SCRIPT_DIR/objectives/$MODULE"
  REPORT="$(go vet -vettool="$(which gost)" ./... 2>&1 || true)"
  echo "---- module $MODULE ----"
  echo "$REPORT"

  SNAPSHOT_PATH="$SCRIPT_DIR/reports/$MODULE.txt"
  if [ "${UPDATE+1}" ]
  then
    if [ ! -z "$REPORT" ]
    then
      echo "$REPORT" > "$SNAPSHOT_PATH"
    else
      if [ -f "$SNAPSHOT_PATH" ]
      then
        rm -- "$SNAPSHOT_PATH"
      fi
    fi
  else
    if [ ! -z "$REPORT" ]
    then
      diff <(echo "$REPORT" | sort) <([ -f "$SNAPSHOT_PATH" ] && cat "$SNAPSHOT_PATH" | sort)
    else
      [ ! -f "$SNAPSHOT_PATH" ]
    fi
  fi
done

