package appserviceenvironments

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type LocalizableString struct {
	LocalizedValue *string `json:"localizedValue,omitempty"`
	Value          *string `json:"value,omitempty"`
}
