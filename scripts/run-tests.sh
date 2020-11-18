#!/bin/sh

go get -d ./... 
go install -v ./...

attempt_counter=0
max_attempts=5

until $(curl --output /dev/null --silent --head --fail -H "Accept: application/vnd.api+json" -XGET http://accountapi:8080/v1/organisation/accounts); do
    if [ ${attempt_counter} -eq ${max_attempts} ]; then
        echo "Max attempts reached"
        exit 1
    fi

    echo 'waiting for accounts API to be up ...'
    attempt_counter=$(($attempt_counter+1))
    sleep 5
done

cd form3
go test -run 'Integration'