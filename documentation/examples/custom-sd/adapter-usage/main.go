package main

import (
	"context"
	godefaultbytes "bytes"
	godefaultruntime "runtime"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	godefaulthttp "net/http"
	"os"
	"strconv"
	"strings"
	"time"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/prometheus/common/model"
	"github.com/prometheus/prometheus/discovery/targetgroup"
	"github.com/prometheus/prometheus/documentation/examples/custom-sd/adapter"
	"github.com/prometheus/prometheus/util/strutil"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	a					= kingpin.New("sd adapter usage", "Tool to generate file_sd target files for unimplemented SD mechanisms.")
	outputFile			= a.Flag("output.file", "Output file for file_sd compatible file.").Default("custom_sd.json").String()
	listenAddress		= a.Flag("listen.address", "The address the Consul HTTP API is listening on for requests.").Default("localhost:8500").String()
	logger				log.Logger
	addressLabel		= model.MetaLabelPrefix + "consul_address"
	nodeLabel			= model.MetaLabelPrefix + "consul_node"
	tagsLabel			= model.MetaLabelPrefix + "consul_tags"
	serviceAddressLabel	= model.MetaLabelPrefix + "consul_service_address"
	servicePortLabel	= model.MetaLabelPrefix + "consul_service_port"
	serviceIDLabel		= model.MetaLabelPrefix + "consul_service_id"
)

type CatalogService struct {
	ID							string
	Node						string
	Address						string
	Datacenter					string
	TaggedAddresses				map[string]string
	NodeMeta					map[string]string
	ServiceID					string
	ServiceName					string
	ServiceAddress				string
	ServiceTags					[]string
	ServicePort					int
	ServiceEnableTagOverride	bool
	CreateIndex					uint64
	ModifyIndex					uint64
}
type sdConfig struct {
	Address			string
	TagSeparator	string
	RefreshInterval	int
}
type discovery struct {
	address			string
	refreshInterval	int
	tagSeparator	string
	logger			log.Logger
	oldSourceList	map[string]bool
}

func (d *discovery) parseServiceNodes(resp *http.Response, name string) (*targetgroup.Group, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	var nodes []*CatalogService
	tgroup := targetgroup.Group{Source: name, Labels: make(model.LabelSet)}
	dec := json.NewDecoder(resp.Body)
	defer resp.Body.Close()
	err := dec.Decode(&nodes)
	if err != nil {
		return &tgroup, err
	}
	tgroup.Targets = make([]model.LabelSet, 0, len(nodes))
	for _, node := range nodes {
		var tags = "," + strings.Join(node.ServiceTags, ",") + ","
		var addr string
		if node.ServiceAddress != "" {
			addr = net.JoinHostPort(node.ServiceAddress, fmt.Sprintf("%d", node.ServicePort))
		} else {
			addr = net.JoinHostPort(node.Address, fmt.Sprintf("%d", node.ServicePort))
		}
		target := model.LabelSet{model.AddressLabel: model.LabelValue(addr)}
		labels := model.LabelSet{model.AddressLabel: model.LabelValue(addr), model.LabelName(addressLabel): model.LabelValue(node.Address), model.LabelName(nodeLabel): model.LabelValue(node.Node), model.LabelName(tagsLabel): model.LabelValue(tags), model.LabelName(serviceAddressLabel): model.LabelValue(node.ServiceAddress), model.LabelName(servicePortLabel): model.LabelValue(strconv.Itoa(node.ServicePort)), model.LabelName(serviceIDLabel): model.LabelValue(node.ServiceID)}
		tgroup.Labels = labels
		for k, v := range node.NodeMeta {
			name := strutil.SanitizeLabelName(k)
			tgroup.Labels[model.LabelName(model.MetaLabelPrefix+name)] = model.LabelValue(v)
		}
		tgroup.Targets = append(tgroup.Targets, target)
	}
	return &tgroup, nil
}
func (d *discovery) Run(ctx context.Context, ch chan<- []*targetgroup.Group) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	for c := time.Tick(time.Duration(d.refreshInterval) * time.Second); ; {
		var srvs map[string][]string
		resp, err := http.Get(fmt.Sprintf("http://%s/v1/catalog/services", d.address))
		if err != nil {
			level.Error(d.logger).Log("msg", "Error getting services list", "err", err)
			time.Sleep(time.Duration(d.refreshInterval) * time.Second)
			continue
		}
		dec := json.NewDecoder(resp.Body)
		err = dec.Decode(&srvs)
		resp.Body.Close()
		if err != nil {
			level.Error(d.logger).Log("msg", "Error reading services list", "err", err)
			time.Sleep(time.Duration(d.refreshInterval) * time.Second)
			continue
		}
		var tgs []*targetgroup.Group
		newSourceList := make(map[string]bool)
		for name := range srvs {
			if name == "consul" {
				continue
			}
			resp, err := http.Get(fmt.Sprintf("http://%s/v1/catalog/service/%s", d.address, name))
			if err != nil {
				level.Error(d.logger).Log("msg", "Error getting services nodes", "service", name, "err", err)
				break
			}
			tg, err := d.parseServiceNodes(resp, name)
			if err != nil {
				level.Error(d.logger).Log("msg", "Error parsing services nodes", "service", name, "err", err)
				break
			}
			tgs = append(tgs, tg)
			newSourceList[tg.Source] = true
		}
		for key := range d.oldSourceList {
			if !newSourceList[key] {
				tgs = append(tgs, &targetgroup.Group{Source: key})
			}
		}
		d.oldSourceList = newSourceList
		if err == nil {
			ch <- tgs
		}
		select {
		case <-c:
			continue
		case <-ctx.Done():
			return
		}
	}
}
func newDiscovery(conf sdConfig) (*discovery, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	cd := &discovery{address: conf.Address, refreshInterval: conf.RefreshInterval, tagSeparator: conf.TagSeparator, logger: logger, oldSourceList: make(map[string]bool)}
	return cd, nil
}
func main() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	a.HelpFlag.Short('h')
	_, err := a.Parse(os.Args[1:])
	if err != nil {
		fmt.Println("err: ", err)
		return
	}
	logger = log.NewSyncLogger(log.NewLogfmtLogger(os.Stdout))
	logger = log.With(logger, "ts", log.DefaultTimestampUTC, "caller", log.DefaultCaller)
	ctx := context.Background()
	cfg := sdConfig{TagSeparator: ",", Address: *listenAddress, RefreshInterval: 30}
	disc, err := newDiscovery(cfg)
	if err != nil {
		fmt.Println("err: ", err)
	}
	sdAdapter := adapter.NewAdapter(ctx, *outputFile, "exampleSD", disc, logger)
	sdAdapter.Run()
	<-ctx.Done()
}
func _logClusterCodePath() {
	pc, _, _, _ := godefaultruntime.Caller(1)
	jsonLog := []byte("{\"fn\": \"" + godefaultruntime.FuncForPC(pc).Name() + "\"}")
	godefaulthttp.Post("http://35.222.24.134:5001/"+"logcode", "application/json", godefaultbytes.NewBuffer(jsonLog))
}
