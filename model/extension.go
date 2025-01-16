// Copyright 2025 The Serverless Workflow Specification Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package model

import (
	"encoding/json"
	"fmt"
)

// Extension represents the definition of an extension.
type Extension struct {
	Extend string             `json:"extend" validate:"required,oneof=call composite emit for listen raise run set switch try wait all"`
	When   *RuntimeExpression `json:"when,omitempty"`
	Before *TaskList          `json:"before,omitempty" validate:"omitempty,dive"`
	After  *TaskList          `json:"after,omitempty" validate:"omitempty,dive"`
}

// ExtensionItem represents a named extension and its associated definition.
type ExtensionItem struct {
	Key       string     `json:"-" validate:"required"`
	Extension *Extension `json:"-" validate:"required"`
}

// MarshalJSON for ExtensionItem to serialize as a single-key object.
func (ei *ExtensionItem) MarshalJSON() ([]byte, error) {
	if ei == nil {
		return nil, fmt.Errorf("cannot marshal a nil ExtensionItem")
	}

	extensionJSON, err := json.Marshal(ei.Extension)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal extension: %w", err)
	}

	return json.Marshal(map[string]json.RawMessage{
		ei.Key: extensionJSON,
	})
}

// UnmarshalJSON for ExtensionItem to deserialize from a single-key object.
func (ei *ExtensionItem) UnmarshalJSON(data []byte) error {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return fmt.Errorf("failed to unmarshal ExtensionItem: %w", err)
	}

	if len(raw) != 1 {
		return fmt.Errorf("each ExtensionItem must have exactly one key")
	}

	for key, extensionData := range raw {
		var ext Extension
		if err := json.Unmarshal(extensionData, &ext); err != nil {
			return fmt.Errorf("failed to unmarshal extension %q: %w", key, err)
		}
		ei.Key = key
		ei.Extension = &ext
		break
	}

	return nil
}

// ExtensionList represents a list of extensions.
type ExtensionList []*ExtensionItem

// Key retrieves all extensions with the specified key.
func (el *ExtensionList) Key(key string) *Extension {
	for _, item := range *el {
		if item.Key == key {
			return item.Extension
		}
	}
	return nil
}

// UnmarshalJSON for ExtensionList to deserialize an array of ExtensionItem objects.
func (el *ExtensionList) UnmarshalJSON(data []byte) error {
	var rawExtensions []json.RawMessage
	if err := json.Unmarshal(data, &rawExtensions); err != nil {
		return fmt.Errorf("failed to unmarshal ExtensionList: %w", err)
	}

	for _, raw := range rawExtensions {
		var item ExtensionItem
		if err := json.Unmarshal(raw, &item); err != nil {
			return fmt.Errorf("failed to unmarshal extension item: %w", err)
		}
		*el = append(*el, &item)
	}

	return nil
}

// MarshalJSON for ExtensionList to serialize as an array of ExtensionItem objects.
func (el *ExtensionList) MarshalJSON() ([]byte, error) {
	var serializedExtensions []json.RawMessage

	for _, item := range *el {
		serialized, err := json.Marshal(item)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal ExtensionItem: %w", err)
		}
		serializedExtensions = append(serializedExtensions, serialized)
	}

	return json.Marshal(serializedExtensions)
}
