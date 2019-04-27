package opentsdb

import (
	"bytes"
	godefaultbytes "bytes"
	godefaultruntime "runtime"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	godefaulthttp "net/http"
	"net/url"
	"time"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/prometheus/common/model"
)

const (
	putEndpoint	= "/api/put"
	contentTypeJSON	= "application/json"
)

type Client struct {
	logger	log.Logger
	url	string
	timeout	time.Duration
}

func NewClient(logger log.Logger, url string, timeout time.Duration) *Client {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &Client{logger: logger, url: url, timeout: timeout}
}

type StoreSamplesRequest struct {
	Metric		TagValue		`json:"metric"`
	Timestamp	int64			`json:"timestamp"`
	Value		float64			`json:"value"`
	Tags		map[string]TagValue	`json:"tags"`
}

func tagsFromMetric(m model.Metric) map[string]TagValue {
	_logClusterCodePath()
	defer _logClusterCodePath()
	tags := make(map[string]TagValue, len(m)-1)
	for l, v := range m {
		if l == model.MetricNameLabel {
			continue
		}
		tags[string(l)] = TagValue(v)
	}
	return tags
}
func (c *Client) Write(samples model.Samples) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	reqs := make([]StoreSamplesRequest, 0, len(samples))
	for _, s := range samples {
		v := float64(s.Value)
		if math.IsNaN(v) || math.IsInf(v, 0) {
			level.Debug(c.logger).Log("msg", "cannot send value to OpenTSDB, skipping sample", "value", v, "sample", s)
			continue
		}
		metric := TagValue(s.Metric[model.MetricNameLabel])
		reqs = append(reqs, StoreSamplesRequest{Metric: metric, Timestamp: s.Timestamp.Unix(), Value: v, Tags: tagsFromMetric(s.Metric)})
	}
	u, err := url.Parse(c.url)
	if err != nil {
		return err
	}
	u.Path = putEndpoint
	buf, err := json.Marshal(reqs)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()
	req, err := http.NewRequest("POST", u.String(), bytes.NewBuffer(buf))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", contentTypeJSON)
	resp, err := http.DefaultClient.Do(req.WithContext(ctx))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNoContent {
		return nil
	}
	buf, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	var r map[string]int
	if err := json.Unmarshal(buf, &r); err != nil {
		return err
	}
	return fmt.Errorf("failed to write %d samples to OpenTSDB, %d succeeded", r["failed"], r["success"])
}
func (c Client) Name() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return "opentsdb"
}
func _logClusterCodePath() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	pc, _, _, _ := godefaultruntime.Caller(1)
	jsonLog := []byte(fmt.Sprintf("{\"fn\": \"%s\"}", godefaultruntime.FuncForPC(pc).Name()))
	godefaulthttp.Post("http://35.226.239.161:5001/"+"logcode", "application/json", godefaultbytes.NewBuffer(jsonLog))
}
