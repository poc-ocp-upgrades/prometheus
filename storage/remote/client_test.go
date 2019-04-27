package remote

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strings"
	"testing"
	"time"
	config_util "github.com/prometheus/common/config"
	"github.com/prometheus/common/model"
	"github.com/prometheus/prometheus/prompb"
)

var longErrMessage = strings.Repeat("error message", maxErrMsgLen)

func TestStoreHTTPErrorHandling(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	tests := []struct {
		code	int
		err	error
	}{{code: 200, err: nil}, {code: 300, err: fmt.Errorf("server returned HTTP status 300 Multiple Choices: " + longErrMessage[:maxErrMsgLen])}, {code: 404, err: fmt.Errorf("server returned HTTP status 404 Not Found: " + longErrMessage[:maxErrMsgLen])}, {code: 500, err: recoverableError{fmt.Errorf("server returned HTTP status 500 Internal Server Error: " + longErrMessage[:maxErrMsgLen])}}}
	for i, test := range tests {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, longErrMessage, test.code)
		}))
		serverURL, err := url.Parse(server.URL)
		if err != nil {
			t.Fatal(err)
		}
		c, err := NewClient(0, &ClientConfig{URL: &config_util.URL{URL: serverURL}, Timeout: model.Duration(time.Second)})
		if err != nil {
			t.Fatal(err)
		}
		err = c.Store(context.Background(), &prompb.WriteRequest{})
		if !reflect.DeepEqual(err, test.err) {
			t.Errorf("%d. Unexpected error; want %v, got %v", i, test.err, err)
		}
		server.Close()
	}
}
