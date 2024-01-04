#!/usr/bin/env bash
set -euo pipefail

readarray -d '' MODULES < <(find . -name go.mod -print0)

RETVAL=0

pushd() {
  command pushd "$@" > /dev/null
}
 
popd() {
  command popd "$@" > /dev/null
}

dir() {
  echo "$(cd "$(dirname "$1")" ; pwd -P)"
}

test() {
  set +e
  for m in "${MODULES[@]}"; do
    pushd "$(dir "$m")"
    go test -count 1 ./... | { grep -v 'no test files'; true; }
    local r=$?
    if [[ $r == 1 ]]; then
      RETVAL=1
    fi
    popd
  done
  set -e
}

tidy() {
  for m in "${MODULES[@]}"; do
    pushd "$(dir "$m")"
    go mod tidy
    local r=$?
    if [[ $r == 1 ]]; then
      RETVAL=1
    fi
    popd
  done
}

while getopts ":h" opt; do
  case ${opt} in
    h)
      echo "Usage:"
      echo "  rw.sh -h    Display this help message"
      echo "  rw.sh test  Test modules"
      echo "  rw.sh tidy  Tidy modules"
      exit 0
      ;;
   \?)
     echo "Invalid Option: -$OPTARG" 1>&2
     exit 1
     ;;
  esac
done
shift $((OPTIND -1))

CMD="$1"; shift
case "$CMD" in
  test)
    test
    ;;
  tidy)
    tidy
    ;;
  *)
    echo "Invalid command $CMD"
    exit 1
    ;;
esac

exit $RETVAL
