#!/bin/bash

FILES=$(find . -type f -name "*.go"|grep -v 'third_party'|grep -v 'testdata')
FAIL=0

for file in  $FILES; do
    error=""
    if [ "$(grep -ic 'copyright' "${file}")" -eq 0 ]; then
        error="        - does not contain a copyright keyword"
        ((FAIL=FAIL+1))
    fi
    if [ "$(grep -i 'copyright' "${file}"|grep SUSE|grep -vc LLC)" -gt 0 ]; then
        error="${error}
        - contains keyword SUSE, but is missing LLC"
        ((FAIL=FAIL+1))
    fi

    if [ "$(grep -ic 'apache' "${file}")" -eq 0 ]; then
        error="${error}
        - has no reference to the apache license"
        ((FAIL=FAIL+1))
    else
        if [ "$(grep -c 'Licensed under the Apache License, Version 2.0 (the "License");' "${file}")" -eq 0 ]; then
            error="${error}
            - is missing the complete the reference to the apache license"
            ((FAIL=FAIL+1))
        fi
        if [ "$(grep -c 'http://www.apache.org/licenses/LICENSE-2.0' "${file}")" -eq 0 ]; then
            error="${error}
            - is missing the link to the apache license"
            ((FAIL=FAIL+1))
        fi
    fi
    if [ -n "${error}" ]; then
    echo "${file}:
${error}
"
    fi
done

if [ "$FAIL" -gt 0 ]; then
    echo "+------------------------------------------------------------------+"
    echo "| Missing or inclomplete copyright/license headers! Please fix it! |"
    echo "+------------------------------------------------------------------+"
    echo " "
    echo " Counted ${FAIL} violations."
    echo " "
    exit 1
fi
echo "+-----------------------------------------------+"
echo "| License & copyright headers seem to be valid. |"
echo "+-----------------------------------------------+"
exit 0
