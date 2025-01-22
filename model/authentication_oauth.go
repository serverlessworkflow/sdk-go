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

// Endpoints are composed here and not on a separate wrapper object to avoid too many nested objects and inline marshaling.
// This allows us to reuse OAuth2AuthenticationProperties also on OpenIdConnectAuthenticationPolicy

type OAuth2AuthenticationPolicy struct {
	Properties *OAuth2AuthenticationProperties `json:",omitempty" validate:"required_without=Use"`
	Endpoints  *OAuth2Endpoints                `json:"endpoints,omitempty"`
	Use        string                          `json:"use,omitempty" validate:"oauth2_policy"`
}

func (o *OAuth2AuthenticationPolicy) ApplyDefaults() {
	if o.Endpoints == nil {
		return
	}

	// Apply defaults if the respective fields are empty
	if o.Endpoints.Token == "" {
		o.Endpoints.Token = OAuth2DefaultTokenURI
	}
	if o.Endpoints.Revocation == "" {
		o.Endpoints.Revocation = OAuth2DefaultRevokeURI
	}
	if o.Endpoints.Introspection == "" {
		o.Endpoints.Introspection = OAuth2DefaultIntrospectionURI
	}
}

func (o *OAuth2AuthenticationPolicy) UnmarshalJSON(data []byte) error {
	type Alias OAuth2AuthenticationPolicy
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(o),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// Initialize Properties if any field for it is set
	if o.Properties == nil && containsOAuth2Properties(data) {
		o.Properties = &OAuth2AuthenticationProperties{}
		if err := json.Unmarshal(data, o.Properties); err != nil {
			return err
		}
	}

	return nil
}

func containsOAuth2Properties(data []byte) bool {
	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return false
	}
	for key := range raw {
		if key != "use" {
			return true
		}
	}
	return false
}

// MarshalJSON customizes the JSON output for OAuth2AuthenticationPolicy
func (o *OAuth2AuthenticationPolicy) MarshalJSON() ([]byte, error) {
	o.ApplyDefaults()
	// Create a map to hold the resulting JSON
	result := make(map[string]interface{})

	// Inline Properties fields if present
	if o.Properties != nil {
		propertiesJSON, err := json.Marshal(o.Properties)
		if err != nil {
			return nil, err
		}

		var propertiesMap map[string]interface{}
		if err := json.Unmarshal(propertiesJSON, &propertiesMap); err != nil {
			return nil, err
		}

		for key, value := range propertiesMap {
			result[key] = value
		}
	}

	// Add the Use field if present
	if o.Use != "" {
		result["use"] = o.Use
	}

	return json.Marshal(result)
}

type OAuth2AuthenticationProperties struct {
	Authority URITemplate                     `json:"authority,omitempty"`
	Grant     OAuth2AuthenticationDataGrant   `json:"grant,omitempty" validate:"oneof='authorization_code' 'client_credentials' 'password' 'refresh_token' 'urn:ietf:params:oauth:grant-type:token-exchange'"`
	Client    *OAuth2AutenthicationDataClient `json:"client,omitempty"`
	Request   *OAuth2TokenRequest             `json:"request,omitempty"`
	Issuers   []string                        `json:"issuers,omitempty"`
	Scopes    []string                        `json:"scopes,omitempty"`
	Audiences []string                        `json:"audiences,omitempty"`
	Username  string                          `json:"username,omitempty"`
	Password  string                          `json:"password,omitempty"`
	Subject   *OAuth2Token                    `json:"subject,omitempty"`
	Actor     *OAuth2Token                    `json:"actor,omitempty"`
}

func (o *OAuth2AuthenticationProperties) UnmarshalJSON(data []byte) error {
	type Alias OAuth2AuthenticationProperties
	aux := &struct {
		Authority json.RawMessage `json:"authority"`
		*Alias
	}{
		Alias: (*Alias)(o),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return fmt.Errorf("failed to unmarshal OAuth2AuthenticationProperties: %w", err)
	}

	// Unmarshal the Authority field
	if aux.Authority != nil {
		uri, err := UnmarshalURITemplate(aux.Authority)
		if err != nil {
			return fmt.Errorf("invalid authority URI: %w", err)
		}
		o.Authority = uri
	}

	return nil
}

// OAuth2AuthenticationDataGrant represents the grant type to use in OAuth2 authentication.
type OAuth2AuthenticationDataGrant string

// Valid grant types
const (
	AuthorizationCodeGrant OAuth2AuthenticationDataGrant = "authorization_code"
	ClientCredentialsGrant OAuth2AuthenticationDataGrant = "client_credentials"
	PasswordGrant          OAuth2AuthenticationDataGrant = "password"
	RefreshTokenGrant      OAuth2AuthenticationDataGrant = "refresh_token"
	TokenExchangeGrant     OAuth2AuthenticationDataGrant = "urn:ietf:params:oauth:grant-type:token-exchange" // #nosec G101
)

type OAuthClientAuthenticationType string

const (
	OAuthClientAuthClientSecretBasic OAuthClientAuthenticationType = "client_secret_basic"
	OAuthClientAuthClientSecretPost  OAuthClientAuthenticationType = "client_secret_post"
	OAuthClientAuthClientSecretJWT   OAuthClientAuthenticationType = "client_secret_jwt"
	OAuthClientAuthPrivateKeyJWT     OAuthClientAuthenticationType = "private_key_jwt"
	OAuthClientAuthNone              OAuthClientAuthenticationType = "none"
)

type OAuth2TokenRequestEncodingType string

const (
	EncodingTypeFormUrlEncoded  OAuth2TokenRequestEncodingType = "application/x-www-form-urlencoded"
	EncodingTypeApplicationJson OAuth2TokenRequestEncodingType = "application/json"
)

// OAuth2AutenthicationDataClient The definition of an OAuth2 client.
type OAuth2AutenthicationDataClient struct {
	ID             string                        `json:"id,omitempty"`
	Secret         string                        `json:"secret,omitempty"`
	Assertion      string                        `json:"assertion,omitempty"`
	Authentication OAuthClientAuthenticationType `json:"authentication,omitempty" validate:"client_auth_type"`
}

type OAuth2TokenRequest struct {
	Encoding OAuth2TokenRequestEncodingType `json:"encoding" validate:"encoding_type"`
}

// OAuth2Token Represents an OAuth2 token.
type OAuth2Token struct {
	// Token The security token to use
	Token string `json:"token,omitempty"`
	// Type The type of the security token to use.
	Type string `json:"type,omitempty"`
}

type OAuth2Endpoints struct {
	Token         string `json:"token,omitempty"`
	Revocation    string `json:"revocation,omitempty"`
	Introspection string `json:"introspection,omitempty"`
}

const (
	OAuth2DefaultTokenURI         = "/oauth2/token" // #nosec G101
	OAuth2DefaultRevokeURI        = "/oauth2/revoke"
	OAuth2DefaultIntrospectionURI = "/oauth2/introspect"
)
