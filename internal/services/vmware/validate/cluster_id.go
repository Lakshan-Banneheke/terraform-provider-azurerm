package validate

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/vmware/2022-05-01/clusters"
)

func ClusterID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := clusters.ParseClusterID(v); err != nil {
		errors = append(errors, err)
	}

	return
}
