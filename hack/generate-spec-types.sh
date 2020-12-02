#!/bin/bash
# Copyright 2020 The Serverless Workflow Specification Authors
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

command -v gojsonschema >/dev/null || go install -modfile=tools.mod -v github.com/atombender/go-jsonschema/cmd/gojsonschema

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
# shellcheck disable=SC2016
sed -i 's/$id/id/g' "${targetdir}/schema/retries.json"


./bin/gojsonschema -v \
  --schema-package=https://serverlessworkflow.org/core/common.json=github.com/serverlessworkflow/sdk-go/model \
  --schema-output=https://serverlessworkflow.org/core/common.json=zz_generated.types_common.go \
  --schema-package=https://serverlessworkflow.org/core/events.json=github.com/serverlessworkflow/sdk-go/model \
  --schema-output=https://serverlessworkflow.org/core/events.json=zz_generated.types_events.go \
  --schema-package=https://serverlessworkflow.org/core/functions.json=github.com/serverlessworkflow/sdk-go/model \
  --schema-output=https://serverlessworkflow.org/core/functions.json=zz_generated.types_functions.go \
  --schema-package=https://serverlessworkflow.org/core/workflow.json=github.com/serverlessworkflow/sdk-go/model \
<<<<<<< HEAD
<<<<<<< HEAD
  --schema-output=https://serverlessworkflow.org/core/workflow.json=zz_generated.types_workflow.go \
  "${targetdir}"/schema/common.json "${targetdir}"/schema/events.json "${targetdir}"/schema/functions.json "${targetdir}"/schema/workflow.json
=======
=======
>>>>>>> configured the generator and added the new generated file for the retries
   --schema-output=https://serverlessworkflow.org/core/workflow.json=zz_generated.types_workflow.go \
  --schema-package=https://serverlessworkflow.org/core/retries.json=github.com/serverlessworkflow/sdk-go/model \
   --schema-output=https://serverlessworkflow.org/core/retries.json=zz_generated.types_retries.go \
  "${targetdir}"/schema/common.json "${targetdir}"/schema/events.json "${targetdir}"/schema/functions.json "${targetdir}"/schema/workflow.json "${targetdir}"/schema/retries.json
<<<<<<< HEAD
>>>>>>> configured the generator and added the new generated file for the retries
=======
>>>>>>> configured the generator and added the new generated file for the retries

sed -i '/type Workflow/d' zz_generated.types_workflow.go

mv -v zz_generated.types_*.go "./${package}/"

cp -v ./hack/zz_generated.types_state_impl.go.template "./${package}/zz_generated.types_state_impl.go"
declare operations=("Delaystate" "Eventstate" "Operationstate" "Parallelstate" "Subflowstate" "Injectstate" "Foreachstate" "Callbackstate" "Databasedswitch" "Eventbasedswitch")
for op in "${operations[@]}"; do
  sed "s/{state}/${op}/g" ./hack/state_interface_impl.template >>"./${package}/zz_generated.types_state_impl.go"
done

go fmt ./...

make addheaders
