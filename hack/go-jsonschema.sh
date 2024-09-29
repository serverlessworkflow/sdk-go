#!/usr/bin/env bash
# Copyright 2022 The Serverless Workflow Specification Authors
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# retrieved from https://github.com/kubernetes/code-generator/blob/master/generate-internal-groups.sh
# and adapted to only install and run the deepcopy-gen

set -o errexit
set -o nounset
set -o pipefail

SCRIPT_ROOT=$(dirname "${BASH_SOURCE[0]}")/..
echo "Script root is $SCRIPT_ROOT"


# To support running this script from anywhere, first cd into this directory,
# and then install with forced module mode on and fully qualified name.
# make sure your GOPATH env is properly set.
# it will go under $GOPATH/bin
cd "$(dirname "${0}")"
GO_JSONSCHEMA_VERSION="v0.16.0"
GO111MODULE=on go install github.com/atombender/go-jsonschema@${GO_JSONSCHEMA_VERSION}

echo "Generating go structs"
mkdir -p ../kubernetes/spec
export GO111MODULE=on
"${GOPATH}/bin/go-jsonschema" \
  -p spec \
  --tags json \
  --struct-name-from-title \
  ../kubernetes/spec/schema.json > ../kubernetes/spec/spec.go
