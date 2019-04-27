package openstack

import (
	"testing"
	"github.com/prometheus/common/model"
	"github.com/prometheus/prometheus/util/testutil"
)

type OpenstackSDHypervisorTestSuite struct{ Mock *SDMock }

func (s *OpenstackSDHypervisorTestSuite) TearDownSuite() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	s.Mock.ShutdownServer()
}
func (s *OpenstackSDHypervisorTestSuite) SetupTest(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	s.Mock = NewSDMock(t)
	s.Mock.Setup()
	s.Mock.HandleHypervisorListSuccessfully()
	s.Mock.HandleVersionsSuccessfully()
	s.Mock.HandleAuthSuccessfully()
}
func (s *OpenstackSDHypervisorTestSuite) openstackAuthSuccess() (Discovery, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	conf := SDConfig{IdentityEndpoint: s.Mock.Endpoint(), Password: "test", Username: "test", DomainName: "12345", Region: "RegionOne", Role: "hypervisor"}
	return NewDiscovery(&conf, nil)
}
func TestOpenstackSDHypervisorRefresh(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	mock := &OpenstackSDHypervisorTestSuite{}
	mock.SetupTest(t)
	hypervisor, _ := mock.openstackAuthSuccess()
	tg, err := hypervisor.refresh()
	testutil.Ok(t, err)
	testutil.Assert(t, tg != nil, "")
	testutil.Assert(t, tg.Targets != nil, "")
	testutil.Assert(t, len(tg.Targets) == 2, "")
	testutil.Equals(t, tg.Targets[0]["__address__"], model.LabelValue("172.16.70.14:0"))
	testutil.Equals(t, tg.Targets[0]["__meta_openstack_hypervisor_hostname"], model.LabelValue("nc14.cloud.com"))
	testutil.Equals(t, tg.Targets[0]["__meta_openstack_hypervisor_type"], model.LabelValue("QEMU"))
	testutil.Equals(t, tg.Targets[0]["__meta_openstack_hypervisor_host_ip"], model.LabelValue("172.16.70.14"))
	testutil.Equals(t, tg.Targets[0]["__meta_openstack_hypervisor_state"], model.LabelValue("up"))
	testutil.Equals(t, tg.Targets[0]["__meta_openstack_hypervisor_status"], model.LabelValue("enabled"))
	testutil.Equals(t, tg.Targets[1]["__address__"], model.LabelValue("172.16.70.13:0"))
	testutil.Equals(t, tg.Targets[1]["__meta_openstack_hypervisor_hostname"], model.LabelValue("cc13.cloud.com"))
	testutil.Equals(t, tg.Targets[1]["__meta_openstack_hypervisor_type"], model.LabelValue("QEMU"))
	testutil.Equals(t, tg.Targets[1]["__meta_openstack_hypervisor_host_ip"], model.LabelValue("172.16.70.13"))
	testutil.Equals(t, tg.Targets[1]["__meta_openstack_hypervisor_state"], model.LabelValue("up"))
	testutil.Equals(t, tg.Targets[1]["__meta_openstack_hypervisor_status"], model.LabelValue("enabled"))
	mock.TearDownSuite()
}
