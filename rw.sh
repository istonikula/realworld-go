#!/usr/bin/env bash
set -euo pipefail

usage() {
  echo "Usage:"
  echo "  rw.sh -h       Display this help message"
  echo "  rw.sh proto    Parse proto files and generate go files"
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

proto_cmd() {
  local protos
  readarray -d '' protos < <(find . -name *.proto -print0)
  for p in "${protos[@]}"; do
    pushd "$(dir "$p")"
    local f="$(basename "$p")"
    protoc \
      --go_out=. --go_opt=paths=source_relative \
      --go_opt=paths=source_relative --go_opt=paths=source_relative \
      "$f"
    popd
  done
}

lint_cmd() {
  echo "lint "$(pwd)""
  golangci-lint run 
}

tidy_cmd() {
  echo "tidy "$(pwd)""
  go mod tidy
}

declare test_cmd_opts="-count=1 -p=1"
test_cmd() {
  go test ${test_cmd_opts} ./... | { grep -v 'no test files'; true; }
}

update_cmd() {
  echo "update and tidy "$(pwd)""
  go get -t -u ./... && go mod tidy
}

run_in_module_dir() {
  local cmd="$1"
  set +e
  for m in "${MODULES[@]}"; do
    pushd "$(dir "$m")"
    $cmd
    if [[ $? == 1 ]]; then
      RETVAL=1
    fi
    popd
  done
  set -e
}

test() {
  while getopts ":v" opt; do
    case ${opt} in
      v)
        test_cmd_opts="${test_cmd_opts} -v"
        ;;
      \?)
        echo "test: Invalid Option: -$OPTARG" 1>&2
        exit 1
        ;;
    esac
  done
  shift $((OPTIND -1))

  run_in_module_dir test_cmd
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
  proto)
    proto_cmd
    ;;
  lint)
    run_in_module_dir lint_cmd
    ;;
  test)
    test "$@"
    ;;
  tidy)
    run_in_module_dir tidy_cmd
    ;;
  update)
    run_in_module_dir update_cmd
    ;;
  *)
    echo "Invalid command $CMD"
    exit 1
    ;;
esac

exit $RETVAL
