#! /bin/bash

CURRENCY_ENDPOINT="http://localhost:${PORT:-8080}/info/currency?"
FAILED=0

test_currency() {
    local params=$1
    local ref=$2

    local endpoint="${CURRENCY_ENDPOINT}${params}"
    echo "Testing endpoint /info/currency?${params}:"

    local api_output=$(curl -s ${CURRENCY_ENDPOINT}${params})

    if [ $? -ne 0 ]; then
    	echo "Failed to process GET request. ret = $?"
    	((FAILED++))
	return
    fi

    diff "${ref}" <(echo "${api_output}") > /dev/null 2>&1

    if [ $? -ne 0 ]; then
	echo "Error: incorrect output on /info/currency?${params} endpoint"
	echo "Expected: $(cat $ref)"
	echo "Got: ${api_output}"
	((FAILED++))
    else
	echo "/info/currency?${params} output: Correct."
    fi
}

for REF_NAME in ./tests/currency/refs/*; do
    QUERRY_PARAMS="${REF_NAME##*/}"

    test_currency "${QUERRY_PARAMS}" "${REF_NAME}"
done

exit ${FAILED}
