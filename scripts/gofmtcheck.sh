#!/usr/bin/env bash

# Check gofmt
files_to_check=$(find . -name '*.go' | grep -v vendor)
gofmt_files=$(${GO_FMT_CMD} -l ${files_to_check})
if [[ -n ${gofmt_files} ]]; then
    echo 'gofmt needs running on the following files:'
    echo "${gofmt_files}"
    echo "You can use the command: \`make fmt\` to reformat code."
    exit 1
fi

exit 0