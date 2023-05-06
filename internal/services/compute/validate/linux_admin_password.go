package validate

import (
	"fmt"
	"regexp"
	"strings"
)

// LinuxAdminPassword validates that admin_password meets the Azure API requirements for Linux Virtual Machines.
func LinuxAdminPassword(i interface{}, k string) (warnings []string, errors []error) {
	// adminPassword must be a string.
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string but it wasn't", k))
		return
	}

	// adminPassword must not be empty.
	if strings.TrimSpace(v) == "" {
		errors = append(errors, fmt.Errorf("%q must not be empty", k))
		return
	}

	// adminPassword Min-length is 6 characters and Max-length is 72 characters.
	if len(v) < 6 || len(v) > 72 {
		errors = append(errors, fmt.Errorf("%q most be between %d and %d characters, got %d", k, 6, 72, len(v)))
	}

	// adminPassword cannot match the following disallowed names.
	disallowedNames := []string{"abc@123", "P@$$w0rd", "P@ssw0rd", "P@ssword123", "Pa$$word", "pass@word1", "Password!", "Password1", "Password22", "iloveyou!"}
	for _, value := range disallowedNames {
		if value == v {
			errors = append(errors, fmt.Errorf("%q specified is not allowed, got %q, cannot match: %q", k, v, strings.Join(disallowedNames, ", ")))
		}
	}

	// adminPassword has to fulfill 3 out of these 4 conditions: Has lower characters, Has upper characters, Has a digit, Has a special character (Regex match [\W_])
	conditions := 0
	tests := []string{"[a-z]", "[A-Z]", "[0-9]", "[^\\d\\w]"}
	for _, test := range tests {
		t, _ := regexp.MatchString(test, v)
		if t {
			conditions++
		}
	}
	if conditions < 3 {
		errors = append(errors, fmt.Errorf("%q has to fulfill 3 out of these 4 conditions: Has lower characters, Has upper characters, Has a digit, Has a special character other than \"_\", fullfiled only %d conditions", k, conditions))
	}

	return warnings, errors
}
