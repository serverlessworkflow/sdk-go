#!/bin/bash
# Copyright 2025 The Serverless Workflow Specification Authors
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


# Script to fetch workflow examples, parse, and validate them using the Go parser.

# Variables
SPEC_REPO="https://github.com/serverlessworkflow/specification"
EXAMPLES_DIR="examples"
PARSER_BINARY="./parser/cmd/main.go"
JUNIT_FILE="./integration-test-junit.xml"

# Create a temporary directory
TEMP_DIR=$(mktemp -d)

# Ensure temporary directory was created
if [ ! -d "$TEMP_DIR" ]; then
    echo "❌ Failed to create a temporary directory."
    exit 1
fi

# shellcheck disable=SC2317
# Clean up the temporary directory on script exit
cleanup() {
    echo "🧹 Cleaning up temporary directory..."
    rm -rf "$TEMP_DIR"
}
trap cleanup EXIT

# Fetch the examples directory
echo "📥 Fetching workflow examples from ${SPEC_REPO}/${EXAMPLES_DIR}..."
if ! git clone --depth=1 --filter=blob:none --sparse "$SPEC_REPO" "$TEMP_DIR" &> /dev/null; then
    echo "❌ Failed to clone specification repository."
    exit 1
fi

cd "$TEMP_DIR" || exit
if ! git sparse-checkout set "$EXAMPLES_DIR" &> /dev/null; then
    echo "❌ Failed to checkout examples directory."
    exit 1
fi

cd - || exit

# Prepare JUnit XML output
echo '<?xml version="1.0" encoding="UTF-8"?>' > "$JUNIT_FILE"
echo '<testsuites>' >> "$JUNIT_FILE"

# Initialize test summary
total_tests=0
failed_tests=0

# Walk through files and validate
echo "⚙️  Running parser on fetched examples..."
while IFS= read -r file; do
    filename=$(basename "$file")
    echo "🔍 Validating: $filename"

    # Run the parser for the file
    if go run "$PARSER_BINARY" "$file" > "$TEMP_DIR/validation.log" 2>&1; then
        echo "✅ Validation succeeded for $filename"
        echo "  <testcase name=\"$filename\" classname=\"integration-test\" />" >> "$JUNIT_FILE"
    else
        echo "❌ Validation failed for $filename"
        failure_message=$(cat "$TEMP_DIR/validation.log" | sed 's/&/&amp;/g; s/</&lt;/g; s/>/&gt;/g')
        echo "  <testcase name=\"$filename\" classname=\"integration-test\">" >> "$JUNIT_FILE"
        echo "    <failure><![CDATA[$failure_message]]></failure>" >> "$JUNIT_FILE"
        echo "  </testcase>" >> "$JUNIT_FILE"
        ((failed_tests++))
    fi

    ((total_tests++))
done < <(find "$TEMP_DIR/$EXAMPLES_DIR" -type f \( -name "*.yaml" -o -name "*.yml" -o -name "*.json" \))

# Finalize JUnit XML output
echo '</testsuites>' >> "$JUNIT_FILE"

# Display test summary
if [ $failed_tests -ne 0 ]; then
    echo "❌ Validation failed for $failed_tests out of $total_tests workflows."
    exit 1
else
    echo "✅ All $total_tests workflows validated successfully."
fi

exit 0
