#!/bin/sh

export GO111MODULE=on
DIR=$(cd $(dirname $0); pwd)
BIN_DIR=$(cd $(dirname $(dirname $0)); pwd)/bin

RESULT=0

mkdir -p ${BIN_DIR}

for absolute in ${DIR}/*/; do
  relative=${absolute#*_example}
  relative=${relative%/}
  go build -o ${BIN_DIR}${relative} ${absolute}/main.go || RESULT=1
done

exit $RESULT
