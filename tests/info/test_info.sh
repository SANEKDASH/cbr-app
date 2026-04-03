#!/bin/bash

FAILED=0
compare_values() {
    local val=$1
    local ref=$2
    local val_name=$3

    if [[ "${val}" != "${ref}" ]]; then
	echo "Error: incorrect ${val_name}: ${val}"
	echo "Expected: ${ref}"
	((FAILED++))
    else
	echo "${val_name}: Correct."
    fi
}


INFO_ENDPOINT="http://localhost:${PORT:-8080}/info"

REF_VERSION=${VERSION:-"1.0.0"}
REF_AUTHOR=${AUTHOR:-"Alex Dash"}
REF_SERVICE="currency"

echo "Processing GET request on ${INFO_ENDPOINT} endpoint"
INFO_RESPONSE=$(curl -s "${INFO_ENDPOINT}")

if [ $? -ne 0 ]; then
    echo "Failed to process GET request."
    exit 1
fi
echo "Success."

INFO_VERSION=$(echo "${INFO_RESPONSE}" | jq -r '.version')
INFO_SERVICE=$(echo "${INFO_RESPONSE}" | jq -r '.service')
INFO_AUTHOR=$(echo "${INFO_RESPONSE}" | jq -r '.author')

compare_values "${INFO_VERSION}" "${REF_VERSION}" "version"
compare_values "${INFO_SERVICE}" "${REF_SERVICE}" "service"
compare_values "${INFO_AUTHOR}" "${REF_AUTHOR}" "author"

exit $FAILED
