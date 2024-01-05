#!/usr/bin/env bash
set -euo pipefail

usage() {
  echo "Usage:"
  echo "  rw.sh -h       Display this help message"
  echo "  rw.sh lint     Lint modules"
  echo "  rw.sh test     Test modules"
  echo "  rw.sh test -v  Test modules, verbose mode"
  echo "  rw.sh tidy     Tidy modules"
  echo "  rw.sh update   Update (and tidy) modules"
}

if [[ $# -lt 1 ]]; then
  usage
  exit 1
fi

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

lint() {
  set +e
  for m in "${MODULES[@]}"; do
    pushd "$(dir "$m")"
    golangci-lint run
    local r=$?
    if [[ $r == 1 ]]; then
      RETVAL=1
    fi
    popd
  done
  set -e
}

test() {
  local verbose=""

  while getopts ":v" opt; do
    case ${opt} in
      v)
        verbose="-v"
        ;;
      \?)
        echo "test: Invalid Option: -$OPTARG" 1>&2
        exit 1
        ;;
    esac
  done
  shift $((OPTIND -1))

  set +e
  for m in "${MODULES[@]}"; do
    pushd "$(dir "$m")"
    go test -count=1 ${verbose} ./... | { grep -v 'no test files'; true; }
    local r=$?
    if [[ $r == 1 ]]; then
      RETVAL=1
    fi
    popd
  done
  set -e
}

tidy() {
  set +e
  for m in "${MODULES[@]}"; do
    pushd "$(dir "$m")"
    go mod tidy
    local r=$?
    if [[ $r == 1 ]]; then
      RETVAL=1
    fi
    popd
  done
  set -e
}

update() {
  set +e
  for m in "${MODULES[@]}"; do
    pushd "$(dir "$m")"
    go get -t -u ./... && go mod tidy
    local r=$?
    if [[ $r == 1 ]]; then
      RETVAL=1
    fi
    popd
  done
  set -e
}

while getopts ":h" opt; do
  case ${opt} in
    h)
      usage
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
  lint)
    lint
    ;;
  test)
    test "$@"
    ;;
  tidy)
    tidy
    ;;
  update)
    update
    ;;
  *)
    echo "Invalid command $CMD"
    exit 1
    ;;
esac

exit $RETVAL
