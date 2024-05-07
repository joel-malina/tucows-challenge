#!/bin/bash

# Copyright 2016 The Kubernetes Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

set -o errexit
set -o nounset
set -o pipefail

export CGO_ENABLED=1
export GO111MODULE=on

GOTESTSUM="gotestsum"

if ! command -v $GOTESTSUM &> /dev/null ; then
    GOTESTSUM="gen/tools/gotestsum"
    if ! command -v $GOTESTSUM &> /dev/null ; then
        go install gotest.tools/gotestsum@latest
    fi
fi


apk add postgresql-client

echo "Testing connection to postgres"
pg_isready -h postgres -p 5432 -U postgres
echo "Running tests:"
$GOTESTSUM --junitfile gen/test/report.xml -- -race -tags musl "$@"
echo
