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
	"strings"
)

// AuthType can be "basic", "bearer", or "oauth2". Default is "basic"
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

// Auth definitions can be used to define authentication information that should be applied to resources
// defined in the operation property of function definitions. It is not used as authentication information
// for the function invocation, but just to access the resource containing the function invocation information.
type Auth struct {
	// Unique auth definition name.
	// +kubebuilder:validation:Required
	Name string `json:"name" validate:"required"`
	// Auth scheme, can be "basic", "bearer", or "oauth2". Default is "basic"
	// +kubebuilder:validation:Enum=basic;bearer;oauth2
	// +kubebuilder:default=basic
	// +kubebuilder:validation:Required
	Scheme AuthType `json:"scheme" validate:"min=1"`
	// Auth scheme properties. Can be one of "Basic properties definition", "Bearer properties definition",
	// or "OAuth2 properties definition"
	// +kubebuilder:validation:Required
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

	switch a.Scheme {
	case AuthTypeBasic:
		authProperties := &BasicAuthProperties{}

		if err := unmarshalKey("properties", auth, authProperties); err != nil {
			return err
		}
		a.Properties.Basic = authProperties

		return nil

	case AuthTypeBearer:
		authProperties := &BearerAuthProperties{}
		if err := unmarshalKey("properties", auth, authProperties); err != nil {
			return err
		}
		a.Properties.Bearer = authProperties
		return nil

	case AuthTypeOAuth2:
		authProperties := &OAuth2AuthProperties{}
		if err := unmarshalKey("properties", auth, authProperties); err != nil {
			return err
		}
		a.Properties.OAuth2 = authProperties
		return nil

	default:
		return fmt.Errorf("failed to parse auth properties")
	}
}

func (a *Auth) MarshalJSON() ([]byte, error) {
	custom, err := json.Marshal(&struct {
		Name       string         `json:"name" validate:"required"`
		Scheme     AuthType       `json:"scheme,omitempty" validate:"omitempty,min=1"`
		Properties AuthProperties `json:"properties" validate:"required"`
	}{
		Name:       a.Name,
		Scheme:     a.Scheme,
		Properties: a.Properties,
	})
	if err != nil {
		fmt.Println(err)
	}
	st := strings.Replace(string(custom), "null,", "", 1)
	st = strings.Replace(st, "\"Basic\":", "", 1)
	st = strings.Replace(st, "\"Oauth2\":", "", 1)
	st = strings.Replace(st, "\"Bearer\":", "", 1)
	st = strings.Replace(st, "{{", "{", 1)
	st = strings.TrimSuffix(st, "}")
	return []byte(st), nil
}

// AuthProperties ...
type AuthProperties struct {
	Basic  *BasicAuthProperties  `json:",omitempty"`
	Bearer *BearerAuthProperties `json:",omitempty"`
	OAuth2 *OAuth2AuthProperties `json:",omitempty"`
}

// BasicAuthProperties Basic Auth Info
type BasicAuthProperties struct {
	Common `json:",inline"`
	// Secret Expression referencing a workflow secret that contains all needed auth info
	// +optional
	Secret string `json:"secret,omitempty"`
	// Username String or a workflow expression. Contains the username
	// +kubebuilder:validation:Required
	Username string `json:"username" validate:"required"`
	// Password String or a workflow expression. Contains the user password
	// +kubebuilder:validation:Required
	Password string `json:"password" validate:"required"`
}

// UnmarshalJSON ...
func (b *BasicAuthProperties) UnmarshalJSON(data []byte) error {
	properties := make(map[string]json.RawMessage)
	if err := json.Unmarshal(data, &properties); err != nil {
		return err
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
	if err := unmarshalKey("secret", properties, &b.Secret); err != nil {
		return err
	}
	return nil
}

// BearerAuthProperties Bearer auth information
type BearerAuthProperties struct {
	Common `json:",inline"`
	// Secret Expression referencing a workflow secret that contains all needed auth info
	// +optional
	Secret string `json:"secret,omitempty"`
	// Token String or a workflow expression. Contains the token
	// +kubebuilder:validation:Required
	Token string `json:"token" validate:"required"`
}

// UnmarshalJSON ...
func (b *BearerAuthProperties) UnmarshalJSON(data []byte) error {
	properties := make(map[string]json.RawMessage)
	if err := json.Unmarshal(data, &properties); err != nil {
		return err
	}
	if err := unmarshalKey("token", properties, &b.Token); err != nil {
		return err
	}
	if err := unmarshalKey("metadata", properties, &b.Metadata); err != nil {
		return err
	}
	if err := unmarshalKey("secret", properties, &b.Secret); err != nil {
		return err
	}
	return nil
}

// OAuth2AuthProperties OAuth2 information
type OAuth2AuthProperties struct {
	Common `json:",inline"`
	// Expression referencing a workflow secret that contains all needed auth info.
	// +optional
	Secret string `json:"secret,omitempty"`
	// String or a workflow expression. Contains the authority information.
	// +optional
	Authority string `json:"authority,omitempty" validate:"omitempty,min=1"`
	// 	Defines the grant type. Can be "password", "clientCredentials", or "tokenExchange"
	// +kubebuilder:validation:Enum=password;clientCredentials;tokenExchange
	// +kubebuilder:validation:Required
	GrantType GrantType `json:"grantType" validate:"required"`
	// String or a workflow expression. Contains the client identifier.
	// +kubebuilder:validation:Required
	ClientID string `json:"clientId" validate:"required"`
	// Workflow secret or a workflow expression. Contains the client secret.
	// +optional
	ClientSecret string `json:"clientSecret,omitempty" validate:"omitempty,min=1"`
	// Array containing strings or workflow expressions. Contains the OAuth2 scopes.
	// +optional
	Scopes []string `json:"scopes,omitempty" validate:"omitempty,min=1"`
	// String or a workflow expression. Contains the username. Used only if grantType is 'resourceOwner'.
	// +optional
	Username string `json:"username,omitempty" validate:"omitempty,min=1"`
	// String or a workflow expression. Contains the user password. Used only if grantType is 'resourceOwner'.
	// +optional
	Password string `json:"password,omitempty" validate:"omitempty,min=1"`
	// Array containing strings or workflow expressions. Contains the OAuth2 audiences.
	// +optional
	Audiences []string `json:"audiences,omitempty" validate:"omitempty,min=1"`
	// String or a workflow expression. Contains the subject token.
	// +optional
	SubjectToken string `json:"subjectToken,omitempty" validate:"omitempty,min=1"`
	// String or a workflow expression. Contains the requested subject.
	// +optional
	RequestedSubject string `json:"requestedSubject,omitempty" validate:"omitempty,min=1"`
	// String or a workflow expression. Contains the requested issuer.
	// +optional
	RequestedIssuer string `json:"requestedIssuer,omitempty" validate:"omitempty,min=1"`
}

// TODO: use reflection to unmarshal the keys and think on a generic approach to handle them

// UnmarshalJSON ...
func (b *OAuth2AuthProperties) UnmarshalJSON(data []byte) error {
	properties := make(map[string]json.RawMessage)
	if err := json.Unmarshal(data, &properties); err != nil {
		return err
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
	if err := unmarshalKey("secret", properties, &b.Secret); err != nil {
		return err
	}
	return nil
}
