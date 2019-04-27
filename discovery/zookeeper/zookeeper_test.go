package zookeeper

import (
	"testing"
	"time"
	"github.com/prometheus/common/model"
)

func TestNewDiscoveryError(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	_, err := NewDiscovery([]string{"unreachable.test"}, time.Second, []string{"/"}, nil, func(data []byte, path string) (model.LabelSet, error) {
		return nil, nil
	})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}
