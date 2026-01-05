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
	"errors"
)

// AuthenticationPolicy Defines an authentication policy.
type AuthenticationPolicy struct {
	Basic       *BasicAuthenticationPolicy         `json:"basic,omitempty"`
	Bearer      *BearerAuthenticationPolicy        `json:"bearer,omitempty"`
	ProxyBearer *ProxyBearerAuthenticationPolicy   `json:"proxy_bearer,omitempty"`
	Digest      *DigestAuthenticationPolicy        `json:"digest,omitempty"`
	OAuth2      *OAuth2AuthenticationPolicy        `json:"oauth2,omitempty"`
	OIDC        *OpenIdConnectAuthenticationPolicy `json:"oidc,omitempty"`
}

// UnmarshalJSON for AuthenticationPolicy to enforce "oneOf" behavior.
func (ap *AuthenticationPolicy) UnmarshalJSON(data []byte) error {
	// Create temporary maps to detect which field is populated
	temp := struct {
		Basic       json.RawMessage `json:"basic"`
		Bearer      json.RawMessage `json:"bearer"`
		ProxyBearer json.RawMessage `json:"proxy_bearer"`
		Digest      json.RawMessage `json:"digest"`
		OAuth2      json.RawMessage `json:"oauth2"`
		OIDC        json.RawMessage `json:"oidc"`
	}{}

	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	// Count non-nil fields
	count := 0
	if len(temp.Basic) > 0 {
		count++
		ap.Basic = &BasicAuthenticationPolicy{}
		if err := json.Unmarshal(temp.Basic, ap.Basic); err != nil {
			return err
		}
	}
	if len(temp.Bearer) > 0 {
		count++
		ap.Bearer = &BearerAuthenticationPolicy{}
		if err := json.Unmarshal(temp.Bearer, ap.Bearer); err != nil {
			return err
		}
	}
	if len(temp.ProxyBearer) > 0 {
		count++
		ap.ProxyBearer = &ProxyBearerAuthenticationPolicy{}
		if err := json.Unmarshal(temp.ProxyBearer, ap.ProxyBearer); err != nil {
			return err
		}
	}
	if len(temp.Digest) > 0 {
		count++
		ap.Digest = &DigestAuthenticationPolicy{}
		if err := json.Unmarshal(temp.Digest, ap.Digest); err != nil {
			return err
		}
	}
	if len(temp.OAuth2) > 0 {
		count++
		ap.OAuth2 = &OAuth2AuthenticationPolicy{}
		if err := json.Unmarshal(temp.OAuth2, ap.OAuth2); err != nil {
			return err
		}
	}
	if len(temp.OIDC) > 0 {
		count++
		ap.OIDC = &OpenIdConnectAuthenticationPolicy{}
		if err := json.Unmarshal(temp.OIDC, ap.OIDC); err != nil {
			return err
		}
	}

	// Ensure only one field is set
	if count != 1 {
		return errors.New("invalid AuthenticationPolicy: only one authentication type must be specified")
	}
	return nil
}

// MarshalJSON for AuthenticationPolicy.
func (ap *AuthenticationPolicy) MarshalJSON() ([]byte, error) {
	if ap.Basic != nil {
		return json.Marshal(map[string]interface{}{"basic": ap.Basic})
	}
	if ap.Bearer != nil {
		return json.Marshal(map[string]interface{}{"bearer": ap.Bearer})
	}
	if ap.ProxyBearer != nil {
		return json.Marshal(map[string]interface{}{"proxy_bearer": ap.ProxyBearer})
	}
	if ap.Digest != nil {
		return json.Marshal(map[string]interface{}{"digest": ap.Digest})
	}
	if ap.OAuth2 != nil {
		return json.Marshal(map[string]interface{}{"oauth2": ap.OAuth2})
	}
	if ap.OIDC != nil {
		return json.Marshal(map[string]interface{}{"oidc": ap.OIDC})
	}
	// Add logic for other fields...
	return nil, errors.New("invalid AuthenticationPolicy: no valid configuration to marshal")
}

// ReferenceableAuthenticationPolicy represents a referenceable authentication policy.
type ReferenceableAuthenticationPolicy struct {
	Use                  *string               `json:"use,omitempty"`
	AuthenticationPolicy *AuthenticationPolicy `json:",inline"`
}

// UnmarshalJSON for ReferenceableAuthenticationPolicy enforces the "oneOf" behavior.
func (rap *ReferenceableAuthenticationPolicy) UnmarshalJSON(data []byte) error {
	// Temporary structure to detect which field is populated
	temp := struct {
		Use *string `json:"use"`
	}{}

	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	// Check if `use` is set
	if temp.Use != nil {
		rap.Use = temp.Use
		return nil
	}

	// If `use` is not set, try unmarshalling inline AuthenticationPolicy
	var ap AuthenticationPolicy
	if err := json.Unmarshal(data, &ap); err != nil {
		return err
	}

	rap.AuthenticationPolicy = &ap
	return nil
}

// MarshalJSON for ReferenceableAuthenticationPolicy.
func (rap *ReferenceableAuthenticationPolicy) MarshalJSON() ([]byte, error) {
	if rap.Use != nil {
		return json.Marshal(map[string]interface{}{"use": rap.Use})
	}
	if rap.AuthenticationPolicy != nil {
		return json.Marshal(rap.AuthenticationPolicy)
	}
	return nil, errors.New("invalid ReferenceableAuthenticationPolicy: no valid configuration to marshal")
}

func NewBasicAuth(username, password string) *AuthenticationPolicy {
	return &AuthenticationPolicy{Basic: &BasicAuthenticationPolicy{
		Username: username,
		Password: password,
	}}
}

// BasicAuthenticationPolicy supports either inline properties (username/password) or a secret reference (use).
type BasicAuthenticationPolicy struct {
	Username string `json:"username,omitempty" validate:"required_without=Use"`
	Password string `json:"password,omitempty" validate:"required_without=Use"`
	Use      string `json:"use,omitempty" validate:"required_without_all=Username Password,basic_policy"`
}

// BearerAuthenticationPolicy supports either an inline token or a secret reference (use).
type BearerAuthenticationPolicy struct {
	Token string `json:"token,omitempty" validate:"required_without=Use,bearer_policy"`
	Use   string `json:"use,omitempty" validate:"required_without=Token"`
}

// ProxyBearerAuthenticationPolicy supports either an inline token or a secret reference (use).
// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Proxy-Authorization
type ProxyBearerAuthenticationPolicy struct {
	Token string `json:"token,omitempty" validate:"required_without=Use,proxy_bearer_policy"`
	Use   string `json:"use,omitempty" validate:"required_without=Token"`
}

// DigestAuthenticationPolicy supports either inline properties (username/password) or a secret reference (use).
type DigestAuthenticationPolicy struct {
	Username string `json:"username,omitempty" validate:"required_without=Use"`
	Password string `json:"password,omitempty" validate:"required_without=Use"`
	Use      string `json:"use,omitempty" validate:"required_without_all=Username Password,digest_policy"`
}

// OpenIdConnectAuthenticationPolicy Use OpenIdConnect authentication.
type OpenIdConnectAuthenticationPolicy struct {
	Properties *OAuth2AuthenticationProperties `json:",omitempty" validate:"omitempty,required_without=Use"`
	Use        string                          `json:"use,omitempty" validate:"omitempty,required_without=Properties"`
}
