#!/usr/bin/env sh

##
## Input parameters
##
BINARY=/stchaind/${BINARY:-stchaind}
ID=${ID:-0}
LOG=${LOG:-stchaind.log}

##
## Assert linux binary
##
if ! [ -f "${BINARY}" ]; then
	echo "The binary $(basename "${BINARY}") cannot be found. Please add the binary to the shared folder. Please use the BINARY environment variable if the name of the binary is not 'stchaind' E.g.: -e BINARY=stchaind_my_test_version"
	exit 1
fi
BINARY_CHECK="$(file "$BINARY" | grep 'ELF 64-bit LSB executable, x86-64')"
if [ -z "${BINARY_CHECK}" ]; then
	echo "Binary needs to be OS linux, ARCH amd64"
	exit 1
fi

##
## Run binary with all parameters
##
export STRATOS_CHAIND_HOME="/stchaind/node${ID}/stchaind"

if [ -d "$(dirname "${STRATOS_CHAIND_HOME}"/"${LOG}")" ]; then
  "${BINARY}" --home "${STRATOS_CHAIND_HOME}" "$@" | tee "${STRATOS_CHAIND_HOME}/${LOG}"
else
  "${BINARY}" --home "${STRATOS_CHAIND_HOME}" "$@"
fi

