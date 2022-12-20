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

GENS="$1"
shift 1

(
  # To support running this script from anywhere, first cd into this directory,
  # and then install with forced module mode on and fully qualified name.
  # make sure your GOPATH env is properly set.
  # it will go under $GOPATH/bin
  cd "$(dirname "${0}")"
  GO111MODULE=on go install k8s.io/code-generator/cmd/deepcopy-gen@latest
)

function codegen::join() { local IFS="$1"; shift; echo "$*"; }

if [ "${GENS}" = "all" ] || grep -qw "deepcopy" <<<"${GENS}"; then
  echo "Generating deepcopy funcs"
  export GO111MODULE=on
  # for debug purposes, increase the log level by updating the -v flag to higher numbers, e.g. -v 4
  "${GOPATH}/bin/deepcopy-gen" -v 1 \
      --input-dirs ./model -O zz_generated.deepcopy \
      --go-header-file "${SCRIPT_ROOT}/hack/boilerplate.txt" \
      "$@"
fi
