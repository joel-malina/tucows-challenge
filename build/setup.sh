#!/bin/bash

GOLINT="golangci-lint"
GOLINT_VERSION=1.54.2

has_right_version() {
    local output=$(command $GOLINT version)
    if [[ "$output" == *"$GOLINT_VERSION"* ]];then
			echo true
		else
			echo false
		fi
}

install() {
    if ! command -v curl &> /dev/null ; then
            wget -O- -nv https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b gen/tools v${GOLINT_VERSION}
        else
            curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b gen/tools v${GOLINT_VERSION}
    fi
}


# Install golang-lint if not present
if ! command -v golangci-lint &> /dev/null ; then
    if [[ ! -f "gen/tools/golangci-lint" ]] ; then
        install
    fi
fi

# Install the specified version of golang-lint
if ! $(has_right_version); then
    install
fi

if ! command -v gotestsum &> /dev/null ; then
    go install gotest.tools/gotestsum@latest
fi

if ! command -v govulncheck &> /dev/null ; then
    go install golang.org/x/vuln/cmd/govulncheck@latest
fi

# go mod download
