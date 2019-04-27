package testutil

import (
	"io/ioutil"
	"os"
	"time"
	"github.com/prometheus/common/model"
	"github.com/prometheus/prometheus/storage"
	"github.com/prometheus/prometheus/storage/tsdb"
)

func NewStorage(t T) storage.Storage {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	dir, err := ioutil.TempDir("", "test_storage")
	if err != nil {
		t.Fatalf("Opening test dir failed: %s", err)
	}
	db, err := tsdb.Open(dir, nil, nil, &tsdb.Options{MinBlockDuration: model.Duration(24 * time.Hour), MaxBlockDuration: model.Duration(24 * time.Hour)})
	if err != nil {
		t.Fatalf("Opening test storage failed: %s", err)
	}
	return testStorage{Storage: tsdb.Adapter(db, int64(0)), dir: dir}
}

type testStorage struct {
	storage.Storage
	dir	string
}

func (s testStorage) Close() error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if err := s.Storage.Close(); err != nil {
		return err
	}
	return os.RemoveAll(s.dir)
}
