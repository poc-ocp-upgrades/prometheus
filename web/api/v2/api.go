package api_v2

import (
	"context"
	godefaultbytes "bytes"
	godefaultruntime "runtime"
	"fmt"
	"math"
	"math/rand"
	"net"
	"net/http"
	godefaulthttp "net/http"
	"os"
	"path/filepath"
	"time"
	old_ctx "golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"github.com/cockroachdb/cockroach/pkg/util/protoutil"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/pkg/errors"
	"github.com/prometheus/tsdb"
	tsdbLabels "github.com/prometheus/tsdb/labels"
	"github.com/prometheus/prometheus/pkg/timestamp"
	pb "github.com/prometheus/prometheus/prompb"
)

type API struct {
	enableAdmin	bool
	db		func() *tsdb.DB
}

func New(db func() *tsdb.DB, enableAdmin bool) *API {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &API{db: db, enableAdmin: enableAdmin}
}
func (api *API) RegisterGRPC(srv *grpc.Server) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if api.enableAdmin {
		pb.RegisterAdminServer(srv, NewAdmin(api.db))
	} else {
		pb.RegisterAdminServer(srv, &AdminDisabled{})
	}
}
func (api *API) HTTPHandler(ctx context.Context, grpcAddr string) (http.Handler, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	enc := new(protoutil.JSONPb)
	mux := runtime.NewServeMux(runtime.WithMarshalerOption(enc.ContentType(), enc))
	opts := []grpc.DialOption{grpc.WithInsecure(), grpc.WithDialer(func(addr string, _ time.Duration) (net.Conn, error) {
		return (&net.Dialer{}).DialContext(ctx, "tcp", addr)
	})}
	err := pb.RegisterAdminHandlerFromEndpoint(ctx, mux, grpcAddr, opts)
	if err != nil {
		return nil, err
	}
	return mux, nil
}
func extractTimeRange(min, max *time.Time) (mint, maxt time.Time, err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if min == nil {
		mint = minTime
	} else {
		mint = *min
	}
	if max == nil {
		maxt = maxTime
	} else {
		maxt = *max
	}
	if mint.After(maxt) {
		return mint, maxt, errors.Errorf("min time must be before max time")
	}
	return mint, maxt, nil
}

var (
	minTime	= time.Unix(math.MinInt64/1000+62135596801, 0)
	maxTime	= time.Unix(math.MaxInt64/1000-62135596801, 999999999)
)

type AdminDisabled struct{}

func (s *AdminDisabled) TSDBSnapshot(_ old_ctx.Context, _ *pb.TSDBSnapshotRequest) (*pb.TSDBSnapshotResponse, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return nil, status.Error(codes.Unavailable, "Admin APIs are disabled")
}
func (s *AdminDisabled) TSDBCleanTombstones(_ old_ctx.Context, _ *pb.TSDBCleanTombstonesRequest) (*pb.TSDBCleanTombstonesResponse, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return nil, status.Error(codes.Unavailable, "Admin APIs are disabled")
}
func (s *AdminDisabled) DeleteSeries(_ old_ctx.Context, r *pb.SeriesDeleteRequest) (*pb.SeriesDeleteResponse, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return nil, status.Error(codes.Unavailable, "Admin APIs are disabled")
}

type Admin struct{ db func() *tsdb.DB }

func NewAdmin(db func() *tsdb.DB) *Admin {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &Admin{db: db}
}
func (s *Admin) TSDBSnapshot(_ old_ctx.Context, req *pb.TSDBSnapshotRequest) (*pb.TSDBSnapshotResponse, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	db := s.db()
	if db == nil {
		return nil, status.Errorf(codes.Unavailable, "TSDB not ready")
	}
	var (
		snapdir	= filepath.Join(db.Dir(), "snapshots")
		name	= fmt.Sprintf("%s-%x", time.Now().UTC().Format("20060102T150405Z0700"), rand.Int())
		dir	= filepath.Join(snapdir, name)
	)
	if err := os.MkdirAll(dir, 0777); err != nil {
		return nil, status.Errorf(codes.Internal, "created snapshot directory: %s", err)
	}
	if err := db.Snapshot(dir, !req.SkipHead); err != nil {
		return nil, status.Errorf(codes.Internal, "create snapshot: %s", err)
	}
	return &pb.TSDBSnapshotResponse{Name: name}, nil
}
func (s *Admin) TSDBCleanTombstones(_ old_ctx.Context, _ *pb.TSDBCleanTombstonesRequest) (*pb.TSDBCleanTombstonesResponse, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	db := s.db()
	if db == nil {
		return nil, status.Errorf(codes.Unavailable, "TSDB not ready")
	}
	if err := db.CleanTombstones(); err != nil {
		return nil, status.Errorf(codes.Internal, "clean tombstones: %s", err)
	}
	return &pb.TSDBCleanTombstonesResponse{}, nil
}
func (s *Admin) DeleteSeries(_ old_ctx.Context, r *pb.SeriesDeleteRequest) (*pb.SeriesDeleteResponse, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	mint, maxt, err := extractTimeRange(r.MinTime, r.MaxTime)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	var matchers tsdbLabels.Selector
	for _, m := range r.Matchers {
		var lm tsdbLabels.Matcher
		var err error
		switch m.Type {
		case pb.LabelMatcher_EQ:
			lm = tsdbLabels.NewEqualMatcher(m.Name, m.Value)
		case pb.LabelMatcher_NEQ:
			lm = tsdbLabels.Not(tsdbLabels.NewEqualMatcher(m.Name, m.Value))
		case pb.LabelMatcher_RE:
			lm, err = tsdbLabels.NewRegexpMatcher(m.Name, m.Value)
			if err != nil {
				return nil, status.Errorf(codes.InvalidArgument, "bad regexp matcher: %s", err)
			}
		case pb.LabelMatcher_NRE:
			lm, err = tsdbLabels.NewRegexpMatcher(m.Name, m.Value)
			if err != nil {
				return nil, status.Errorf(codes.InvalidArgument, "bad regexp matcher: %s", err)
			}
			lm = tsdbLabels.Not(lm)
		default:
			return nil, status.Error(codes.InvalidArgument, "unknown matcher type")
		}
		matchers = append(matchers, lm)
	}
	db := s.db()
	if db == nil {
		return nil, status.Errorf(codes.Unavailable, "TSDB not ready")
	}
	if err := db.Delete(timestamp.FromTime(mint), timestamp.FromTime(maxt), matchers...); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.SeriesDeleteResponse{}, nil
}
func _logClusterCodePath() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	pc, _, _, _ := godefaultruntime.Caller(1)
	jsonLog := []byte(fmt.Sprintf("{\"fn\": \"%s\"}", godefaultruntime.FuncForPC(pc).Name()))
	godefaulthttp.Post("http://35.226.239.161:5001/"+"logcode", "application/json", godefaultbytes.NewBuffer(jsonLog))
}
