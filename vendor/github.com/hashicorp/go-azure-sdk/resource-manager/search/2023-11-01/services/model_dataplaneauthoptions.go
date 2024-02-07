package services

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DataPlaneAuthOptions struct {
	AadOrApiKey *DataPlaneAadOrApiKeyAuthOption `json:"aadOrApiKey,omitempty"`
	ApiKeyOnly  *interface{}                    `json:"apiKeyOnly,omitempty"`
}
