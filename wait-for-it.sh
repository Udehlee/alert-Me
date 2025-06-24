#!/bin/sh

WAITFORIT_cmdname=${0##*/}

echoerr() {
  if [ "$WAITFORIT_QUIET" -ne 1 ]; then
    echo "$@" 1>&2
  fi
}

usage() {
  cat << USAGE >&2
Usage:
    $WAITFORIT_cmdname host:port [-s] [-t timeout] [-- command args]
    -h HOST | --host=HOST       Host or IP under test
    -p PORT | --port=PORT       TCP port under test
    -s | --strict               Only execute subcommand if the test succeeds
    -q | --quiet                Don't output any status messages
    -t TIMEOUT | --timeout=TIMEOUT Timeout in seconds, zero for no timeout
    -- COMMAND ARGS             Execute command with args after the test finishes
USAGE
  exit 1
}

wait_for() {
  if [ "$WAITFORIT_TIMEOUT" -gt 0 ]; then
    echoerr "$WAITFORIT_cmdname: waiting $WAITFORIT_TIMEOUT seconds for $WAITFORIT_HOST:$WAITFORIT_PORT"
  else
    echoerr "$WAITFORIT_cmdname: waiting for $WAITFORIT_HOST:$WAITFORIT_PORT without a timeout"
  fi
  start_ts=$(date +%s)
  while :
  do
    nc -z "$WAITFORIT_HOST" "$WAITFORIT_PORT" > /dev/null 2>&1
    result=$?
    if [ "$result" -eq 0 ]; then
      end_ts=$(date +%s)
      echoerr "$WAITFORIT_cmdname: $WAITFORIT_HOST:$WAITFORIT_PORT is available after $((end_ts - start_ts)) seconds"
      break
    fi
    sleep 1
  done
  return $result
}

# Parse args
WAITFORIT_TIMEOUT=15
WAITFORIT_STRICT=0
WAITFORIT_QUIET=0
WAITFORIT_CLI=""
while [ $# -gt 0 ]; do
  case "$1" in
    *:* )
      WAITFORIT_HOST=$(echo "$1" | cut -d: -f1)
      WAITFORIT_PORT=$(echo "$1" | cut -d: -f2)
      shift
      ;;
    -q|--quiet)
      WAITFORIT_QUIET=1
      shift
      ;;
    -s|--strict)
      WAITFORIT_STRICT=1
      shift
      ;;
    -h)
      WAITFORIT_HOST="$2"
      shift 2
      ;;
    --host=*)
      WAITFORIT_HOST="${1#*=}"
      shift
      ;;
    -p)
      WAITFORIT_PORT="$2"
      shift 2
      ;;
    --port=*)
      WAITFORIT_PORT="${1#*=}"
      shift
      ;;
    -t)
      WAITFORIT_TIMEOUT="$2"
      shift 2
      ;;
    --timeout=*)
      WAITFORIT_TIMEOUT="${1#*=}"
      shift
      ;;
    --)
      shift
      WAITFORIT_CLI="$@"
      break
      ;;
    *)
      echoerr "Unknown argument: $1"
      usage
      ;;
  esac
done

if [ -z "$WAITFORIT_HOST" ] || [ -z "$WAITFORIT_PORT" ]; then
  echoerr "Error: you need to provide a host and port to test."
  usage
fi

wait_for
WAITFORIT_RESULT=$?

if [ "$WAITFORIT_RESULT" -ne 0 ] && [ "$WAITFORIT_STRICT" -eq 1 ]; then
  echoerr "$WAITFORIT_cmdname: strict mode, refusing to execute subprocess"
  exit $WAITFORIT_RESULT
fi

if [ -n "$WAITFORIT_CLI" ]; then
  exec $WAITFORIT_CLI
else
  exit $WAITFORIT_RESULT
fi
