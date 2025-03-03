package util

import "sigs.k8s.io/controller-runtime/pkg/client"

func MatchesLabels(obj client.Object, labels map[string]string) bool {
	if labels == nil || len(labels) == 0 {
		return true
	}
	for k, v := range labels {
		vv, ok := obj.GetLabels()[k]
		if !ok {
			return false
		}
		if v != vv {
			return false
		}
	}
	return true
}
