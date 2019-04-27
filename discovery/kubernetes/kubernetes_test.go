package kubernetes

import (
	"context"
	"encoding/json"
	"sync"
	"testing"
	"time"
	"github.com/go-kit/kit/log"
	"github.com/prometheus/prometheus/discovery/targetgroup"
	"github.com/prometheus/prometheus/util/testutil"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"
	k8stesting "k8s.io/client-go/testing"
	"k8s.io/client-go/tools/cache"
)

type watcherFactory struct {
	sync.RWMutex
	watchers	map[schema.GroupVersionResource]*watch.FakeWatcher
}

func (wf *watcherFactory) watchFor(gvr schema.GroupVersionResource) *watch.FakeWatcher {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	wf.Lock()
	defer wf.Unlock()
	var fakewatch *watch.FakeWatcher
	fakewatch, ok := wf.watchers[gvr]
	if !ok {
		fakewatch = watch.NewFakeWithChanSize(128, true)
		wf.watchers[gvr] = fakewatch
	}
	return fakewatch
}
func (wf *watcherFactory) Nodes() *watch.FakeWatcher {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return wf.watchFor(schema.GroupVersionResource{Group: "", Version: "v1", Resource: "nodes"})
}
func (wf *watcherFactory) Ingresses() *watch.FakeWatcher {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return wf.watchFor(schema.GroupVersionResource{Group: "extensions", Version: "v1beta1", Resource: "ingresses"})
}
func (wf *watcherFactory) Endpoints() *watch.FakeWatcher {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return wf.watchFor(schema.GroupVersionResource{Group: "", Version: "v1", Resource: "endpoints"})
}
func (wf *watcherFactory) Services() *watch.FakeWatcher {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return wf.watchFor(schema.GroupVersionResource{Group: "", Version: "v1", Resource: "services"})
}
func (wf *watcherFactory) Pods() *watch.FakeWatcher {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return wf.watchFor(schema.GroupVersionResource{Group: "", Version: "v1", Resource: "pods"})
}
func makeDiscovery(role Role, nsDiscovery NamespaceDiscovery, objects ...runtime.Object) (*Discovery, kubernetes.Interface, *watcherFactory) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	clientset := fake.NewSimpleClientset(objects...)
	wf := &watcherFactory{watchers: make(map[schema.GroupVersionResource]*watch.FakeWatcher)}
	clientset.PrependWatchReactor("*", func(action k8stesting.Action) (handled bool, ret watch.Interface, err error) {
		gvr := action.GetResource()
		return true, wf.watchFor(gvr), nil
	})
	return &Discovery{client: clientset, logger: log.NewNopLogger(), role: role, namespaceDiscovery: &nsDiscovery}, clientset, wf
}

type k8sDiscoveryTest struct {
	discovery		discoverer
	beforeRun		func()
	afterStart		func()
	expectedMaxItems	int
	expectedRes		map[string]*targetgroup.Group
}

func (d k8sDiscoveryTest) Run(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	ch := make(chan []*targetgroup.Group)
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	if d.beforeRun != nil {
		d.beforeRun()
	}
	go d.discovery.Run(ctx, ch)
	resChan := make(chan map[string]*targetgroup.Group)
	go readResultWithTimeout(t, ch, d.expectedMaxItems, time.Second, resChan)
	dd, ok := d.discovery.(hasSynced)
	if !ok {
		t.Errorf("discoverer does not implement hasSynced interface")
		return
	}
	if !cache.WaitForCacheSync(ctx.Done(), dd.hasSynced) {
		t.Errorf("discoverer failed to sync: %v", dd)
		return
	}
	if d.afterStart != nil {
		d.afterStart()
	}
	if d.expectedRes != nil {
		res := <-resChan
		requireTargetGroups(t, d.expectedRes, res)
	}
}
func readResultWithTimeout(t *testing.T, ch <-chan []*targetgroup.Group, max int, timeout time.Duration, resChan chan<- map[string]*targetgroup.Group) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	allTgs := make([][]*targetgroup.Group, 0)
Loop:
	for {
		select {
		case tgs := <-ch:
			allTgs = append(allTgs, tgs)
			if len(allTgs) == max {
				break Loop
			}
		case <-time.After(timeout):
			t.Logf("timed out, got %d (max: %d) items, some events are skipped", len(allTgs), max)
			break Loop
		}
	}
	res := make(map[string]*targetgroup.Group)
	for _, tgs := range allTgs {
		for _, tg := range tgs {
			if tg == nil {
				continue
			}
			res[tg.Source] = tg
		}
	}
	resChan <- res
}
func requireTargetGroups(t *testing.T, expected, res map[string]*targetgroup.Group) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	b1, err := json.Marshal(expected)
	if err != nil {
		panic(err)
	}
	b2, err := json.Marshal(res)
	if err != nil {
		panic(err)
	}
	testutil.Equals(t, string(b1), string(b2))
}

type hasSynced interface{ hasSynced() bool }

var _ hasSynced = &Discovery{}
var _ hasSynced = &Node{}
var _ hasSynced = &Endpoints{}
var _ hasSynced = &Ingress{}
var _ hasSynced = &Pod{}
var _ hasSynced = &Service{}

func (d *Discovery) hasSynced() bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	d.RLock()
	defer d.RUnlock()
	for _, discoverer := range d.discoverers {
		if hasSynceddiscoverer, ok := discoverer.(hasSynced); ok {
			if !hasSynceddiscoverer.hasSynced() {
				return false
			}
		}
	}
	return true
}
func (n *Node) hasSynced() bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return n.informer.HasSynced()
}
func (e *Endpoints) hasSynced() bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return e.endpointsInf.HasSynced() && e.serviceInf.HasSynced() && e.podInf.HasSynced()
}
func (i *Ingress) hasSynced() bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return i.informer.HasSynced()
}
func (p *Pod) hasSynced() bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return p.informer.HasSynced()
}
func (s *Service) hasSynced() bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return s.informer.HasSynced()
}
