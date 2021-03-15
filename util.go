package docker

import (
	"fmt"

	"github.com/docker/docker/api/types/filters"
)

// GetSessionLabels gets labels set for all resources created in this session
func (c Client) GetSessionLabels() map[string]string {
	return map[string]string{
		"sessionId":  c.SessionID,
		"created-by": "docker-utils",
	}
}

// FilterByLabels creates filters.Args based on a set of labels
func FilterByLabels(labels map[string]string) (args filters.Args) {
	var kvps []filters.KeyValuePair
	for key, value := range labels {
		kvps = append(kvps, filters.KeyValuePair{
			Key:   "label",
			Value: fmt.Sprintf("%s=%s", key, value),
		})
	}
	return filters.NewArgs(kvps...)
}
