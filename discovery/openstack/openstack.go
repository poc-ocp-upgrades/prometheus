package openstack

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"
	"github.com/go-kit/kit/log"
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/mwitkow/go-conntrack"
	"github.com/prometheus/client_golang/prometheus"
	config_util "github.com/prometheus/common/config"
	"github.com/prometheus/common/model"
	"github.com/prometheus/prometheus/discovery/targetgroup"
)

var (
	refreshFailuresCount	= prometheus.NewCounter(prometheus.CounterOpts{Name: "prometheus_sd_openstack_refresh_failures_total", Help: "The number of OpenStack-SD scrape failures."})
	refreshDuration		= prometheus.NewSummary(prometheus.SummaryOpts{Name: "prometheus_sd_openstack_refresh_duration_seconds", Help: "The duration of an OpenStack-SD refresh in seconds."})
	DefaultSDConfig		= SDConfig{Port: 80, RefreshInterval: model.Duration(60 * time.Second)}
)

type SDConfig struct {
	IdentityEndpoint		string			`yaml:"identity_endpoint"`
	Username			string			`yaml:"username"`
	UserID				string			`yaml:"userid"`
	Password			config_util.Secret	`yaml:"password"`
	ProjectName			string			`yaml:"project_name"`
	ProjectID			string			`yaml:"project_id"`
	DomainName			string			`yaml:"domain_name"`
	DomainID			string			`yaml:"domain_id"`
	ApplicationCredentialName	string			`yaml:"application_credential_name"`
	ApplicationCredentialID		string			`yaml:"application_credential_id"`
	ApplicationCredentialSecret	config_util.Secret	`yaml:"application_credential_secret"`
	Role				Role			`yaml:"role"`
	Region				string			`yaml:"region"`
	RefreshInterval			model.Duration		`yaml:"refresh_interval,omitempty"`
	Port				int			`yaml:"port"`
	AllTenants			bool			`yaml:"all_tenants,omitempty"`
	TLSConfig			config_util.TLSConfig	`yaml:"tls_config,omitempty"`
}
type Role string

const (
	OpenStackRoleHypervisor	Role	= "hypervisor"
	OpenStackRoleInstance	Role	= "instance"
)

func (c *Role) UnmarshalYAML(unmarshal func(interface{}) error) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if err := unmarshal((*string)(c)); err != nil {
		return err
	}
	switch *c {
	case OpenStackRoleHypervisor, OpenStackRoleInstance:
		return nil
	default:
		return fmt.Errorf("unknown OpenStack SD role %q", *c)
	}
}
func (c *SDConfig) UnmarshalYAML(unmarshal func(interface{}) error) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	*c = DefaultSDConfig
	type plain SDConfig
	err := unmarshal((*plain)(c))
	if err != nil {
		return err
	}
	if c.Role == "" {
		return fmt.Errorf("role missing (one of: instance, hypervisor)")
	}
	if c.Region == "" {
		return fmt.Errorf("openstack SD configuration requires a region")
	}
	return nil
}
func init() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	prometheus.MustRegister(refreshFailuresCount)
	prometheus.MustRegister(refreshDuration)
}

type Discovery interface {
	Run(ctx context.Context, ch chan<- []*targetgroup.Group)
	refresh() (tg *targetgroup.Group, err error)
}

func NewDiscovery(conf *SDConfig, l log.Logger) (Discovery, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	var opts gophercloud.AuthOptions
	if conf.IdentityEndpoint == "" {
		var err error
		opts, err = openstack.AuthOptionsFromEnv()
		if err != nil {
			return nil, err
		}
	} else {
		opts = gophercloud.AuthOptions{IdentityEndpoint: conf.IdentityEndpoint, Username: conf.Username, UserID: conf.UserID, Password: string(conf.Password), TenantName: conf.ProjectName, TenantID: conf.ProjectID, DomainName: conf.DomainName, DomainID: conf.DomainID, ApplicationCredentialID: conf.ApplicationCredentialID, ApplicationCredentialName: conf.ApplicationCredentialName, ApplicationCredentialSecret: string(conf.ApplicationCredentialSecret)}
	}
	client, err := openstack.NewClient(opts.IdentityEndpoint)
	if err != nil {
		return nil, err
	}
	tls, err := config_util.NewTLSConfig(&conf.TLSConfig)
	if err != nil {
		return nil, err
	}
	client.HTTPClient = http.Client{Transport: &http.Transport{IdleConnTimeout: 5 * time.Duration(conf.RefreshInterval), TLSClientConfig: tls, DialContext: conntrack.NewDialContextFunc(conntrack.DialWithTracing(), conntrack.DialWithName("openstack_sd"))}, Timeout: 5 * time.Duration(conf.RefreshInterval)}
	switch conf.Role {
	case OpenStackRoleHypervisor:
		hypervisor := NewHypervisorDiscovery(client, &opts, time.Duration(conf.RefreshInterval), conf.Port, conf.Region, l)
		return hypervisor, nil
	case OpenStackRoleInstance:
		instance := NewInstanceDiscovery(client, &opts, time.Duration(conf.RefreshInterval), conf.Port, conf.Region, conf.AllTenants, l)
		return instance, nil
	default:
		return nil, errors.New("unknown OpenStack discovery role")
	}
}
