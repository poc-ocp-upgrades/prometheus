package remote

import (
	"github.com/prometheus/common/model"
	"github.com/prometheus/prometheus/pkg/labels"
	"github.com/prometheus/prometheus/storage"
)

func (s *Storage) Appender() (storage.Appender, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return s, nil
}
func (s *Storage) Add(l labels.Labels, t int64, v float64) (uint64, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	s.mtx.RLock()
	defer s.mtx.RUnlock()
	for _, q := range s.queues {
		if err := q.Append(&model.Sample{Metric: labelsToMetric(l), Timestamp: model.Time(t), Value: model.SampleValue(v)}); err != nil {
			panic(err)
		}
	}
	return 0, nil
}
func (s *Storage) AddFast(l labels.Labels, _ uint64, t int64, v float64) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_, err := s.Add(l, t, v)
	return err
}
func (*Storage) Commit() error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return nil
}
func (*Storage) Rollback() error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return nil
}
