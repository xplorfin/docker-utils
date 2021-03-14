package docker

import (
	"fmt"

	"github.com/docker/docker/api/types/filters"
)

// getSessionLabels gets labels set for all resources created in this session
func (c Client) getSessionLabels() map[string]string {
	return map[string]string{
		"sessionId":  c.SessionID,
		"created-by": "docker-utils",
	}
}

// filterByLabels creates filters.Args based on a set of labels
func filterByLabels(labels map[string]string) (args filters.Args) {
	var kvps []filters.KeyValuePair
	for key, value := range labels {
		kvps = append(kvps, filters.KeyValuePair{
			Key:   "label",
			Value: fmt.Sprintf("%s=%s", key, value),
		})
	}
	return filters.NewArgs(kvps...)
}
