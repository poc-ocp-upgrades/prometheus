package azure

import (
	"reflect"
	"testing"
	"github.com/Azure/azure-sdk-for-go/arm/compute"
)

func TestMapFromVMWithEmptyTags(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	id := "test"
	name := "name"
	vmType := "type"
	location := "westeurope"
	networkProfile := compute.NetworkProfile{}
	properties := &compute.VirtualMachineProperties{StorageProfile: &compute.StorageProfile{OsDisk: &compute.OSDisk{OsType: "Linux"}}, NetworkProfile: &networkProfile}
	testVM := compute.VirtualMachine{ID: &id, Name: &name, Type: &vmType, Location: &location, Tags: nil, Properties: properties}
	expectedVM := virtualMachine{ID: id, Name: name, Type: vmType, Location: location, OsType: "Linux", Tags: map[string]*string{}, NetworkProfile: networkProfile}
	actualVM := mapFromVM(testVM)
	if !reflect.DeepEqual(expectedVM, actualVM) {
		t.Errorf("Expected %v got %v", expectedVM, actualVM)
	}
}
func TestMapFromVMWithTags(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	id := "test"
	name := "name"
	vmType := "type"
	location := "westeurope"
	tags := map[string]*string{"prometheus": new(string)}
	networkProfile := compute.NetworkProfile{}
	properties := &compute.VirtualMachineProperties{StorageProfile: &compute.StorageProfile{OsDisk: &compute.OSDisk{OsType: "Linux"}}, NetworkProfile: &networkProfile}
	testVM := compute.VirtualMachine{ID: &id, Name: &name, Type: &vmType, Location: &location, Tags: &tags, Properties: properties}
	expectedVM := virtualMachine{ID: id, Name: name, Type: vmType, Location: location, OsType: "Linux", Tags: tags, NetworkProfile: networkProfile}
	actualVM := mapFromVM(testVM)
	if !reflect.DeepEqual(expectedVM, actualVM) {
		t.Errorf("Expected %v got %v", expectedVM, actualVM)
	}
}
func TestMapFromVMScaleSetVMWithEmptyTags(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	id := "test"
	name := "name"
	vmType := "type"
	location := "westeurope"
	networkProfile := compute.NetworkProfile{}
	properties := &compute.VirtualMachineScaleSetVMProperties{StorageProfile: &compute.StorageProfile{OsDisk: &compute.OSDisk{OsType: "Linux"}}, NetworkProfile: &networkProfile}
	testVM := compute.VirtualMachineScaleSetVM{ID: &id, Name: &name, Type: &vmType, Location: &location, Tags: nil, Properties: properties}
	scaleSet := "testSet"
	expectedVM := virtualMachine{ID: id, Name: name, Type: vmType, Location: location, OsType: "Linux", Tags: map[string]*string{}, NetworkProfile: networkProfile, ScaleSet: scaleSet}
	actualVM := mapFromVMScaleSetVM(testVM, scaleSet)
	if !reflect.DeepEqual(expectedVM, actualVM) {
		t.Errorf("Expected %v got %v", expectedVM, actualVM)
	}
}
func TestMapFromVMScaleSetVMWithTags(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	id := "test"
	name := "name"
	vmType := "type"
	location := "westeurope"
	tags := map[string]*string{"prometheus": new(string)}
	networkProfile := compute.NetworkProfile{}
	properties := &compute.VirtualMachineScaleSetVMProperties{StorageProfile: &compute.StorageProfile{OsDisk: &compute.OSDisk{OsType: "Linux"}}, NetworkProfile: &networkProfile}
	testVM := compute.VirtualMachineScaleSetVM{ID: &id, Name: &name, Type: &vmType, Location: &location, Tags: &tags, Properties: properties}
	scaleSet := "testSet"
	expectedVM := virtualMachine{ID: id, Name: name, Type: vmType, Location: location, OsType: "Linux", Tags: tags, NetworkProfile: networkProfile, ScaleSet: scaleSet}
	actualVM := mapFromVMScaleSetVM(testVM, scaleSet)
	if !reflect.DeepEqual(expectedVM, actualVM) {
		t.Errorf("Expected %v got %v", expectedVM, actualVM)
	}
}
