package connection

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ConnectionUpdateParameters struct {
	Name       *string                     `json:"name,omitempty"`
	Properties *ConnectionUpdateProperties `json:"properties,omitempty"`
}
