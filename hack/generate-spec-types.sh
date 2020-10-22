#!/bin/bash
# Copyright 2020 The Serverless Workflow Specification Authors
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


command -v ./bin/gojsonschema >/dev/null || go build -o ./bin/gojsonschema github.com/atombender/go-jsonschema/cmd/gojsonschema && go mod tidy

echo "--> Generating specification types"

declare package="model"
declare targetdir="/tmp/serverlessworkflow"

if [ ! -d "${targetdir}" ]; then
  git clone git@github.com:serverlessworkflow/specification.git ${targetdir}
fi

# remove once we have https://github.com/atombender/go-jsonschema/pull/16
# shellcheck disable=SC2016
sed -i 's/$id/id/g' "${targetdir}/schema/common.json"
# shellcheck disable=SC2016
sed -i 's/$id/id/g' "${targetdir}/schema/events.json"
# shellcheck disable=SC2016
sed -i 's/$id/id/g' "${targetdir}/schema/functions.json"
# shellcheck disable=SC2016
sed -i 's/$id/id/g' "${targetdir}/schema/workflow.json"

./bin/gojsonschema -v \
  --schema-package=https://serverlessworkflow.org/core/common.json=github.com/serverlessworkflow/sdk-go/model \
   --schema-output=https://serverlessworkflow.org/core/common.json=generated_types_common.go \
  --schema-package=https://serverlessworkflow.org/core/events.json=github.com/serverlessworkflow/sdk-go/model \
   --schema-output=https://serverlessworkflow.org/core/events.json=generated_types_events.go \
  --schema-package=https://serverlessworkflow.org/core/functions.json=github.com/serverlessworkflow/sdk-go/model \
   --schema-output=https://serverlessworkflow.org/core/functions.json=generated_types_functions.go \
  --schema-package=https://serverlessworkflow.org/core/workflow.json=github.com/serverlessworkflow/sdk-go/model \
   --schema-output=https://serverlessworkflow.org/core/workflow.json=generated_types_workflow.go \
  "${targetdir}"/schema/common.json "${targetdir}"/schema/events.json "${targetdir}"/schema/functions.json "${targetdir}"/schema/workflow.json

sed -i '/type Workflow/d' generated_types_workflow.go

mv -v generated_types_*.go "./${package}/"

go fmt ./...

make addheaders