// Copyright 2021 The Serverless Workflow Specification Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
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

// AuthType ...
type AuthType string

const (
	// AuthTypeBasic ...
	AuthTypeBasic AuthType = "basic"
	// AuthTypeBearer ...
	AuthTypeBearer AuthType = "bearer"
	// AuthTypeOAuth2 ...
	AuthTypeOAuth2 AuthType = "oauth2"
)

// GrantType ...
type GrantType string

const (
	// GrantTypePassword ...
	GrantTypePassword GrantType = "password"
	// GrantTypeClientCredentials ...
	GrantTypeClientCredentials GrantType = "clientCredentials"
	// GrantTypeTokenExchange ...
	GrantTypeTokenExchange GrantType = "tokenExchange"
)

func getAuthProperties(authType AuthType) (AuthProperties, bool) {
	switch authType {
	case AuthTypeBasic:
		return &BasicAuthProperties{}, true
	case AuthTypeBearer:
		return &BearerAuthProperties{}, true
	case AuthTypeOAuth2:
		return &OAuth2AuthProperties{}, true
	}
	return nil, false
}

// Auth ...
type Auth struct {
	// Name Unique auth definition name
	Name string `json:"name" validate:"required"`
	// Scheme Defines the auth type
	Scheme AuthType `json:"scheme,omitempty" validate:"omitempty,min=1"`
	// Properties ...
	Properties AuthProperties `json:"properties" validate:"required"`
}

// UnmarshalJSON Auth definition
func (a *Auth) UnmarshalJSON(data []byte) error {
	auth := make(map[string]json.RawMessage)
	if err := json.Unmarshal(data, &auth); err != nil {
		// it's a file
		file, err := unmarshalFile(data)
		if err != nil {
			return err
		}
		// call us recursively
		if err := json.Unmarshal(file, &a); err != nil {
			return err
		}
		return nil
	}
	if err := unmarshalKey("scheme", auth, &a.Scheme); err != nil {
		return err
	}
	if err := unmarshalKey("name", auth, &a.Name); err != nil {
		return err
	}

	if len(a.Scheme) == 0 {
		a.Scheme = AuthTypeBasic
	}
	authProperties, ok := getAuthProperties(a.Scheme)
	if !ok {
		return fmt.Errorf("authentication scheme %s not supported", a.Scheme)
	}

	// we take the type we want to unmarshal based on the scheme
	if err := unmarshalKey("properties", auth, authProperties); err != nil {
		return err
	}

	a.Properties = authProperties
	return nil
}

// AuthProperties ...
type AuthProperties interface {
	// GetMetadata ...
	GetMetadata() *Metadata
	// GetSecret ...
	GetSecret() string
	// DeepCopyAuthProperties fixes in.Properties.DeepCopyAuthProperties undefined (type AuthProperties has no
	// field or method DeepCopyAuthProperties)
	DeepCopyAuthProperties() AuthProperties
}

// BaseAuthProperties ...
type BaseAuthProperties struct {
	Common
	// Secret Expression referencing a workflow secret that contains all needed auth info
	Secret string `json:"secret,omitempty"`
}

// UnmarshalJSON ...
func (b *BaseAuthProperties) UnmarshalJSON(data []byte) error {
	properties := make(map[string]json.RawMessage)
	if err := json.Unmarshal(data, &properties); err != nil {
		b.Secret, err = unmarshalString(data)
		if err != nil {
			return err
		}
		return nil
	}
	if err := unmarshalKey("metadata", properties, &b.Metadata); err != nil {
		return err
	}
	if err := unmarshalKey("secret", properties, &b.Secret); err != nil {
		return err
	}
	return nil
}

// GetMetadata ...
func (b *BaseAuthProperties) GetMetadata() *Metadata {
	return &b.Metadata
}

// GetSecret ...
func (b *BaseAuthProperties) GetSecret() string {
	return b.Secret
}

func (b *BasicAuthProperties) DeepCopyAuthProperties() AuthProperties {
	return b
}

// BasicAuthProperties Basic Auth Info
type BasicAuthProperties struct {
	BaseAuthProperties
	// Username String or a workflow expression. Contains the username
	Username string `json:"username" validate:"required"`
	// Password String or a workflow expression. Contains the user password
	Password string `json:"password" validate:"required"`
}

// UnmarshalJSON ...
func (b *BasicAuthProperties) UnmarshalJSON(data []byte) error {
	properties := make(map[string]json.RawMessage)
	if err := json.Unmarshal(data, &properties); err != nil {
		err = json.Unmarshal(data, &b.BaseAuthProperties)
		if err != nil {
			return err
		}
		return nil
	}
	if err := unmarshalKey("username", properties, &b.Username); err != nil {
		return err
	}
	if err := unmarshalKey("password", properties, &b.Password); err != nil {
		return err
	}
	if err := unmarshalKey("metadata", properties, &b.Metadata); err != nil {
		return err
	}
	return nil
}

// BearerAuthProperties Bearer auth information
type BearerAuthProperties struct {
	BaseAuthProperties
	// Token String or a workflow expression. Contains the token
	Token string `json:"token" validate:"required"`
}

func (b *BearerAuthProperties) DeepCopyAuthProperties() AuthProperties {
	return b
}

// UnmarshalJSON ...
func (b *BearerAuthProperties) UnmarshalJSON(data []byte) error {
	properties := make(map[string]json.RawMessage)
	if err := json.Unmarshal(data, &properties); err != nil {
		err = json.Unmarshal(data, &b.BaseAuthProperties)
		if err != nil {
			return err
		}
		return nil
	}
	if err := unmarshalKey("token", properties, &b.Token); err != nil {
		return err
	}
	if err := unmarshalKey("metadata", properties, &b.Metadata); err != nil {
		return err
	}
	return nil
}

// OAuth2AuthProperties OAuth2 information
type OAuth2AuthProperties struct {
	BaseAuthProperties
	// Authority String or a workflow expression. Contains the authority information
	Authority string `json:"authority,omitempty" validate:"omitempty,min=1"`
	// GrantType Defines the grant type
	GrantType GrantType `json:"grantType" validate:"required"`
	// ClientID String or a workflow expression. Contains the client identifier
	ClientID string `json:"clientId" validate:"required"`
	// ClientSecret Workflow secret or a workflow expression. Contains the client secret
	ClientSecret string `json:"clientSecret,omitempty" validate:"omitempty,min=1"`
	// Scopes Array containing strings or workflow expressions. Contains the OAuth2 scopes
	Scopes []string `json:"scopes,omitempty" validate:"omitempty,min=1"`
	// Username String or a workflow expression. Contains the username. Used only if grantType is 'resourceOwner'
	Username string `json:"username,omitempty" validate:"omitempty,min=1"`
	// Password String or a workflow expression. Contains the user password. Used only if grantType is 'resourceOwner'
	Password string `json:"password,omitempty" validate:"omitempty,min=1"`
	// Audiences Array containing strings or workflow expressions. Contains the OAuth2 audiences
	Audiences []string `json:"audiences,omitempty" validate:"omitempty,min=1"`
	// SubjectToken String or a workflow expression. Contains the subject token
	SubjectToken string `json:"subjectToken,omitempty" validate:"omitempty,min=1"`
	// RequestedSubject String or a workflow expression. Contains the requested subject
	RequestedSubject string `json:"requestedSubject,omitempty" validate:"omitempty,min=1"`
	// RequestedIssuer String or a workflow expression. Contains the requested issuer
	RequestedIssuer string `json:"requestedIssuer,omitempty" validate:"omitempty,min=1"`
}

func (b *OAuth2AuthProperties) DeepCopyAuthProperties() AuthProperties {
	return b
}

// TODO: use reflection to unmarshal the keys and think on a generic approach to handle them

// UnmarshalJSON ...
func (b *OAuth2AuthProperties) UnmarshalJSON(data []byte) error {
	properties := make(map[string]json.RawMessage)
	if err := json.Unmarshal(data, &properties); err != nil {
		err = json.Unmarshal(data, &b.BaseAuthProperties)
		if err != nil {
			return err
		}
		return nil
	}
	if err := unmarshalKey("authority", properties, &b.Authority); err != nil {
		return err
	}
	if err := unmarshalKey("grantType", properties, &b.GrantType); err != nil {
		return err
	}
	if err := unmarshalKey("clientId", properties, &b.ClientID); err != nil {
		return err
	}
	if err := unmarshalKey("clientSecret", properties, &b.ClientSecret); err != nil {
		return err
	}
	if err := unmarshalKey("scopes", properties, &b.Scopes); err != nil {
		return err
	}
	if err := unmarshalKey("username", properties, &b.Username); err != nil {
		return err
	}
	if err := unmarshalKey("password", properties, &b.Password); err != nil {
		return err
	}
	if err := unmarshalKey("audiences", properties, &b.Audiences); err != nil {
		return err
	}
	if err := unmarshalKey("subjectToken", properties, &b.SubjectToken); err != nil {
		return err
	}
	if err := unmarshalKey("requestedSubject", properties, &b.RequestedSubject); err != nil {
		return err
	}
	if err := unmarshalKey("requestedIssuer", properties, &b.RequestedIssuer); err != nil {
		return err
	}
	if err := unmarshalKey("metadata", properties, &b.Metadata); err != nil {
		return err
	}
	return nil
}
