#!/bin/bash

A24API_BINARY=""
A24API_CONFIG=""
A24API_DOMAIN=""
A24API_TEST=""

SLEEPTIME=2

for CMD_PARAM in "$@"; do
    CMD_PARAM1=${CMD_PARAM%%=*}
    if [ "$CMD_PARAM1" == "--config" -o "$CMD_PARAM1" == "-c" ]; then A24API_CONFIG="${CMD_PARAM#*=}";
    elif [ "$CMD_PARAM1" == "--binary" -o "$CMD_PARAM1" == "-b" ]; then A24API_BINARY="${CMD_PARAM#*=}";
    elif [ "$CMD_PARAM1" == "--domain" -o "$CMD_PARAM1" == "-d" ]; then A24API_DOMAIN="${CMD_PARAM#*=}";
    elif [ "$CMD_PARAM1" == "--test" -o "$CMD_PARAM1" == "-t" ]; then A24API_TEST="${CMD_PARAM#*=}";
    fi
done

if [ "$A24API_BINARY" == "" -o "$A24API_DOMAIN" == "" ]; then
    echo "Usage: $0 [options]

Options:
    -c|--config=<path>              Path to config file.
    -b|--binary=<path>              Path to a24api binary.
    -d|--domain=<name>              Domain for testing.
    -t|--test=<name>                Specify test.
"
    exit 1
fi

if [ "A24API_CONFIG" != "" ]; then export A24API_CONFIG=$A24API_CONFIG; fi

function testList () {
    # list domains
    echo "## list domains"
    $A24API_BINARY dns list
    sleep ${SLEEPTIME}
    # list records
    echo "## list records"
    $A24API_BINARY dns list ${A24API_DOMAIN}
    sleep ${SLEEPTIME}
    # list records - wrong domain name
    echo "## list records wrong domain name"
    $A24API_BINARY dns list ${A24API_DOMAIN}-a24test
    sleep ${SLEEPTIME}
}

function testA () {
    # create A record
    echo "## create A record"
    $A24API_BINARY dns create ${A24API_DOMAIN} A a24test-A 3600 127.0.0.1
    sleep ${SLEEPTIME}
    # create A record - that already exists
    echo "## create A record that already exists"
    $A24API_BINARY dns create ${A24API_DOMAIN} A a24test-A 3600 127.0.0.1
    sleep ${SLEEPTIME}
    # list A record
    echo "## list previously created A record"
    $A24API_BINARY dns list ${A24API_DOMAIN} -ft '^A$' -fn 'a24test-'
    sleep ${SLEEPTIME}
    # update A record
    echo "## update previously created A record"
    $A24API_BINARY dns update $($A24API_BINARY dns list ${A24API_DOMAIN} -ft '^A$' -fn 'a24test-A' | tr -s [:space:] | cut -d" " -f1,2) A a24test-B 1800 127.0.0.2
    sleep ${SLEEPTIME}
    # list A record
    echo "## list previously updated A record"
    $A24API_BINARY dns list ${A24API_DOMAIN} -ft '^A$' -fn 'a24test-'
    sleep ${SLEEPTIME}
    # delete A record
    echo "## delete previously created A record"
    $A24API_BINARY dns delete $($A24API_BINARY dns list ${A24API_DOMAIN} -ft '^A$' -fn 'a24test-' | tr -s [:space:] | cut -d" " -f1,2)
    sleep ${SLEEPTIME}
}

function testAAAA () {
    # create AAAA record
    echo "## create AAAA record"
    $A24API_BINARY dns create ${A24API_DOMAIN} AAAA a24test-A 3600 0:0:0:0:0:0:0:1
    sleep ${SLEEPTIME}
    # create AAAA record - that already exists
    echo "## create AAAA record that already exists"
    $A24API_BINARY dns create ${A24API_DOMAIN} AAAA a24test-A 3600 0:0:0:0:0:0:0:1
    sleep ${SLEEPTIME}
    # list AAAA record
    echo "## list previously created AAAA record"
    $A24API_BINARY dns list ${A24API_DOMAIN} -ft '^AAAA$' -fn 'a24test-'
    sleep ${SLEEPTIME}
    # update AAAA record
    echo "## update previously created AAAA record"
    $A24API_BINARY dns update $($A24API_BINARY dns list ${A24API_DOMAIN} -ft '^AAAA$' -fn 'a24test-' | tr -s [:space:] | cut -d" " -f1,2) AAAA a24test-B 1800 0:0:0:0:0:0:0:2
    sleep ${SLEEPTIME}
    # list AAAA record
    echo "## list previously updated AAAA record"
    $A24API_BINARY dns list ${A24API_DOMAIN} -ft '^AAAA$' -fn 'a24test-'
    sleep ${SLEEPTIME}
    # delete AAAA record
    echo "## delete previously created AAAA record"
    $A24API_BINARY dns delete $($A24API_BINARY dns list ${A24API_DOMAIN} -ft '^AAAA$' -fn 'a24test-' | tr -s [:space:] | cut -d" " -f1,2)
    sleep ${SLEEPTIME}
}

function testCNAME () {
    # create CNAME record
    echo "## create CNAME record"
    $A24API_BINARY dns create ${A24API_DOMAIN} CNAME a24test-A 3600 example.com.
    sleep ${SLEEPTIME}
    # create CNAME record - that already exists
    echo "## create CNAME record that already exists"
    $A24API_BINARY dns create ${A24API_DOMAIN} CNAME a24test-A 3600 example.com.
    sleep ${SLEEPTIME}
    # list CNAME record
    echo "## list previously created CNAME record"
    $A24API_BINARY dns list ${A24API_DOMAIN} -ft '^CNAME$' -fn 'a24test-'
    sleep ${SLEEPTIME}
    # update CNAME record
    echo "## update previously created CNAME record"
    $A24API_BINARY dns update $($A24API_BINARY dns list ${A24API_DOMAIN} -ft '^CNAME$' -fn 'a24test-' | tr -s [:space:] | cut -d" " -f1,2) CNAME a24test-B 1800 example2.com.
    sleep ${SLEEPTIME}
    # list CNAME record
    echo "## list previously updated CNAME record"
    $A24API_BINARY dns list ${A24API_DOMAIN} -ft '^CNAME$' -fn 'a24test-'
    sleep ${SLEEPTIME}
    # delete CNAME record
    echo "## delete previously created CNAME record"
    $A24API_BINARY dns delete $($A24API_BINARY dns list ${A24API_DOMAIN} -ft '^CNAME$' -fn 'a24test-' | tr -s [:space:] | cut -d" " -f1,2)
    sleep ${SLEEPTIME}
}

function testTXT () {
    # create TXT record
    echo "## create TXT record"
    $A24API_BINARY dns create ${A24API_DOMAIN} TXT a24test-A 3600 "Active24 API Test 1"
    sleep ${SLEEPTIME}
    # create TXT record - that already exists
    echo "## create TXT record that already exists"
    $A24API_BINARY dns create ${A24API_DOMAIN} TXT a24test-A 3600 "Active24 API Test 1"
    sleep ${SLEEPTIME}
    # list TXT record
    echo "## list previously created TXT record"
    $A24API_BINARY dns list ${A24API_DOMAIN} -ft '^TXT$' -fn 'a24test-'
    sleep ${SLEEPTIME}
    # update TXT record
    echo "## update previously created TXT record"
    $A24API_BINARY dns update $($A24API_BINARY dns list ${A24API_DOMAIN} -ft '^TXT$' -fn 'a24test-' | tr -s [:space:] | cut -d" " -f1,2) TXT a24test-B 1800 "Active24 API Test 2"
    sleep ${SLEEPTIME}
    # list TXT record
    echo "## list previously updated TXT record"
    $A24API_BINARY dns list ${A24API_DOMAIN} -ft '^TXT$' -fn 'a24test-'
    sleep ${SLEEPTIME}
    # delete TXT record
    echo "## delete previously created TXT record"
    $A24API_BINARY dns delete $($A24API_BINARY dns list ${A24API_DOMAIN} -ft '^TXT$' -fn 'a24test-' | tr -s [:space:] | cut -d" " -f1,2)
    sleep ${SLEEPTIME}
}

function testNS () {
    # create NS record
    echo "## create NS record"
    $A24API_BINARY dns create ${A24API_DOMAIN} NS a24test-A 3600 example.com.
    sleep ${SLEEPTIME}
    # create NS record - that already exists
    echo "## create NS record that already exists"
    $A24API_BINARY dns create ${A24API_DOMAIN} NS a24test-A 3600 example.com.
    sleep ${SLEEPTIME}
    # list NS record
    echo "## list previously created NS record"
    $A24API_BINARY dns list ${A24API_DOMAIN} -ft '^NS$' -fn 'a24test-'
    sleep ${SLEEPTIME}
    # update NS record
    echo "## update previously created NS record"
    $A24API_BINARY dns update $($A24API_BINARY dns list ${A24API_DOMAIN} -ft '^NS$' -fn 'a24test-' | tr -s [:space:] | cut -d" " -f1,2) NS a24test-B 1800 example2.com.
    sleep ${SLEEPTIME}
    # list NS record
    echo "## list previously updated NS record"
    $A24API_BINARY dns list ${A24API_DOMAIN} -ft '^NS$' -fn 'a24test-'
    sleep ${SLEEPTIME}
    # delete NS record
    echo "## delete previously created NS record"
    $A24API_BINARY dns delete $($A24API_BINARY dns list ${A24API_DOMAIN} -ft '^NS$' -fn 'a24test-' | tr -s [:space:] | cut -d" " -f1,2)
    sleep ${SLEEPTIME}
}

function testSSHFP () {
    # create SSHFP record
    echo "## create SSHFP record"
    $A24API_BINARY dns create ${A24API_DOMAIN} SSHFP a24test-A 3600 4 2 4ddf47cf93bf5237974bf27ff363030abe9032a96cc8eed4d877851720d3a11e
    sleep ${SLEEPTIME}
    # create SSHFP record - that already exists
    echo "## create SSHFP record that already exists"
    $A24API_BINARY dns create ${A24API_DOMAIN} SSHFP a24test-A 3600 4 2 4ddf47cf93bf5237974bf27ff363030abe9032a96cc8eed4d877851720d3a11e
    sleep ${SLEEPTIME}
    # list SSHFP record
    echo "## list previously created SSHFP record"
    $A24API_BINARY dns list ${A24API_DOMAIN} -ft '^SSHFP$' -fn 'a24test-'
    sleep ${SLEEPTIME}
    # update SSHFP record
    echo "## update previously created SSHFP record"
    $A24API_BINARY dns update $($A24API_BINARY dns list ${A24API_DOMAIN} -ft '^SSHFP$' -fn 'a24test-' | tr -s [:space:] | cut -d" " -f1,2) SSHFP a24test-B 1800 1 2 088be59412563e97f75a4e1dfa8b98d17eb7b98accabdd8ebbc0f6e8e2bb919d
    sleep ${SLEEPTIME}
    # list SSHFP record
    echo "## list previously updated SSHFP record"
    $A24API_BINARY dns list ${A24API_DOMAIN} -ft '^SSHFP$' -fn 'a24test-'
    sleep ${SLEEPTIME}
    # delete SSHFP record
    echo "## delete previously created SSHFP record"
    $A24API_BINARY dns delete $($A24API_BINARY dns list ${A24API_DOMAIN} -ft '^SSHFP$' -fn 'a24test-' | tr -s [:space:] | cut -d" " -f1,2)
    sleep ${SLEEPTIME}
}

function testSRV () {
    # create SRV record
    echo "## create SRV record"
    $A24API_BINARY dns create ${A24API_DOMAIN} SRV _a24test-A._tcp 3600 10 20 1234 example.com.
    sleep ${SLEEPTIME}
    # create SRV record - that already exists
    echo "## create SRV record that already exists"
    $A24API_BINARY dns create ${A24API_DOMAIN} SRV _a24test-A._tcp 3600 10 20 1234 example.com.
    sleep ${SLEEPTIME}
    # list SRV record
    echo "## list previously created SRV record"
    $A24API_BINARY dns list ${A24API_DOMAIN} -ft '^SRV$' -fn 'a24test-'
    sleep ${SLEEPTIME}
    # update SRV record
    echo "## update previously created SRV record"
    $A24API_BINARY dns update $($A24API_BINARY dns list ${A24API_DOMAIN} -ft '^SRV$' -fn 'a24test-' | tr -s [:space:] | cut -d" " -f1,2) SRV _a24test-B._tcp 1800 20 30 5678 example2.com.
    sleep ${SLEEPTIME}
    # list SRV record
    echo "## list previously updated SRV record"
    $A24API_BINARY dns list ${A24API_DOMAIN} -ft '^SRV$' -fn 'a24test-'
    sleep ${SLEEPTIME}
    # delete SRV record
    echo "## delete previously created SRV record"
    $A24API_BINARY dns delete $($A24API_BINARY dns list ${A24API_DOMAIN} -ft '^SRV$' -fn 'a24test-' | tr -s [:space:] | cut -d" " -f1,2)
    sleep ${SLEEPTIME}
}

function testTLSA () {
    # create TLSA record
    echo "## create TLSA record"
    $A24API_BINARY dns create ${A24API_DOMAIN} TLSA _25._tcp.a24test-A 3600 3 1 1 0C72AC70B745AC19998811B131D662C9AC69DBDBE7CB23E5B514B56664C5D3D6
    sleep ${SLEEPTIME}
    # create TLSA record - that already exists
    echo "## create TLSA record that already exists"
    $A24API_BINARY dns create ${A24API_DOMAIN} TLSA _25._tcp.a24test-A 3600 3 1 1 0C72AC70B745AC19998811B131D662C9AC69DBDBE7CB23E5B514B56664C5D3D6
    sleep ${SLEEPTIME}
    # list TLSA record
    echo "## list previously created TLSA record"
    $A24API_BINARY dns list ${A24API_DOMAIN} -ft '^TLSA$' -fn 'a24test-'
    sleep ${SLEEPTIME}
    # update TLSA record
    echo "## update previously created TLSA record"
    $A24API_BINARY dns update $($A24API_BINARY dns list ${A24API_DOMAIN} -ft '^TLSA$' -fn 'a24test-' | tr -s [:space:] | cut -d" " -f1,2) TLSA _25._tcp.a24test-B 1800 3 0 1 AB9BEB9919729F3239AF08214C1EF6CCA52D2DBAE788BB5BE834C13911292ED9
    sleep ${SLEEPTIME}
    # list TLSA record
    echo "## list previously updated TLSA record"
    $A24API_BINARY dns list ${A24API_DOMAIN} -ft '^TLSA$' -fn 'a24test-'
    sleep ${SLEEPTIME}
    # delete TLSA record
    echo "## delete previously created TLSA record"
    $A24API_BINARY dns delete $($A24API_BINARY dns list ${A24API_DOMAIN} -ft '^TLSA$' -fn 'a24test-' | tr -s [:space:] | cut -d" " -f1,2)
    sleep ${SLEEPTIME}
}

function testCAA () {
    # create CAA record
    echo "## create CAA record"
    $A24API_BINARY dns create ${A24API_DOMAIN} CAA a24test-A 3600 0 issue "ca.example.com"
    sleep ${SLEEPTIME}
    # create CAA record - that already exists
    echo "## create CAA record that already exists"
    $A24API_BINARY dns create ${A24API_DOMAIN} CAA a24test-A 3600 0 issue "ca.example.com"
    sleep ${SLEEPTIME}
    # list CAA record
    echo "## list previously created CAA record"
    $A24API_BINARY dns list ${A24API_DOMAIN} -ft '^CAA$' -fn 'a24test-'
    sleep ${SLEEPTIME}
    # update CAA record
    echo "## update previously created CAA record"
    $A24API_BINARY dns update $($A24API_BINARY dns list ${A24API_DOMAIN} -ft '^CAA$' -fn 'a24test-' | tr -s [:space:] | cut -d" " -f1,2) CAA a24test-B 1800 0 issue "ca.example2.com"
    sleep ${SLEEPTIME}
    # list CAA record
    echo "## list previously updated CAA record"
    $A24API_BINARY dns list ${A24API_DOMAIN} -ft '^CAA$' -fn 'a24test-'
    sleep ${SLEEPTIME}
    # delete CAA record
    echo "## delete previously created CAA record"
    $A24API_BINARY dns delete $($A24API_BINARY dns list ${A24API_DOMAIN} -ft '^CAA$' -fn 'a24test-' | tr -s [:space:] | cut -d" " -f1,2)
    sleep ${SLEEPTIME}
}

function testMX () {
    # create MX record
    echo "## create MX record"
    $A24API_BINARY dns create ${A24API_DOMAIN} MX a24test-A 3600 10 example.com.
    sleep ${SLEEPTIME}
    # create MX record - that already exists
    echo "## create MX record that already exists"
    $A24API_BINARY dns create ${A24API_DOMAIN} MX a24test-A 3600 10 example.com.
    sleep ${SLEEPTIME}
    # list MX record
    echo "## list previously created MX record"
    $A24API_BINARY dns list ${A24API_DOMAIN} -ft '^MX$' -fn 'a24test-'
    sleep ${SLEEPTIME}
    # update MX record
    echo "## update previously created MX record"
    $A24API_BINARY dns update $($A24API_BINARY dns list ${A24API_DOMAIN} -ft '^MX$' -fn 'a24test-' | tr -s [:space:] | cut -d" " -f1,2) MX a24test-B 1800 20 example2.com.
    sleep ${SLEEPTIME}
    # list MX record
    echo "## list previously updated MX record"
    $A24API_BINARY dns list ${A24API_DOMAIN} -ft '^MX$' -fn 'a24test-'
    sleep ${SLEEPTIME}
    # delete MX record
    echo "## delete previously created MX record"
    $A24API_BINARY dns delete $($A24API_BINARY dns list ${A24API_DOMAIN} -ft '^MX$' -fn 'a24test-' | tr -s [:space:] | cut -d" " -f1,2)
    sleep ${SLEEPTIME}
}

if [ "$A24API_TEST" != "" ]; then
    test${A24API_TEST}
else
    testList
    testA
    testAAAA
    testCNAME
    testTXT
    testNS
    testSSHFP
    testSRV
    testTLSA
    testCAA
fi
