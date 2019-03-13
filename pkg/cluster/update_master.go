package cluster

import (
	"bytes"
	"context"
	"sort"
	"strings"

	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2018-10-01/compute"

	"github.com/openshift/openshift-azure/pkg/api"
	"github.com/openshift/openshift-azure/pkg/cluster/updateblob"
	"github.com/openshift/openshift-azure/pkg/config"
)

func (u *simpleUpgrader) filterOldVMs(vms []compute.VirtualMachineScaleSetVM, blob *updateblob.UpdateBlob, ssHash []byte) []compute.VirtualMachineScaleSetVM {
	var oldVMs []compute.VirtualMachineScaleSetVM
	for _, vm := range vms {
		if !bytes.Equal(blob.HostnameHashes[*vm.Name], ssHash) {
			oldVMs = append(oldVMs, vm)
		} else {
			u.log.Infof("skipping vm %q since it's already updated", *vm.Name)
		}
	}
	return oldVMs
}

// UpdateMasterAgentPool updates one by one all the VMs of the master scale set,
// in place.
func (u *simpleUpgrader) UpdateMasterAgentPool(ctx context.Context, cs *api.OpenShiftManagedCluster, app *api.AgentPoolProfile) *api.PluginError {
	ssName := config.MasterScalesetName
	desiredHash, err := u.hasher.HashScaleSet(cs, app)
	if err != nil {
		return &api.PluginError{Err: err, Step: api.PluginStepUpdateMasterAgentPoolHashScaleSet}
	}

	blob, err := u.updateBlobService.Read()
	if err != nil {
		return &api.PluginError{Err: err, Step: api.PluginStepUpdateMasterAgentPoolReadBlob}
	}

	vms, err := u.vmc.List(ctx, cs.Properties.AzProfile.ResourceGroup, ssName, "", "", "")
	if err != nil {
		return &api.PluginError{Err: err, Step: api.PluginStepUpdateMasterAgentPoolListVMs}
	}

	// range our vms in order, so that if we previously crashed half-way through
	// updating one and it is broken, we pick up where we left off.
	sort.Slice(vms, func(i, j int) bool {
		return *vms[i].VirtualMachineScaleSetVMProperties.OsProfile.ComputerName <
			*vms[j].VirtualMachineScaleSetVMProperties.OsProfile.ComputerName
	})

	vms = u.filterOldVMs(vms, blob, desiredHash)
	for _, vm := range vms {
		hostname := strings.ToLower(*vm.VirtualMachineScaleSetVMProperties.OsProfile.ComputerName)
		u.log.Infof("draining %s", hostname)
		err = u.kubeclient.DeleteMaster(hostname)
		if err != nil {
			return &api.PluginError{Err: err, Step: api.PluginStepUpdateMasterAgentPoolDrain}
		}

		u.log.Infof("deallocating %s", hostname)
		err = u.vmc.Deallocate(ctx, cs.Properties.AzProfile.ResourceGroup, ssName, *vm.InstanceID)
		if err != nil {
			return &api.PluginError{Err: err, Step: api.PluginStepUpdateMasterAgentPoolDeallocate}
		}

		u.log.Infof("updating %s", hostname)
		err = u.ssc.UpdateInstances(ctx, cs.Properties.AzProfile.ResourceGroup, ssName, compute.VirtualMachineScaleSetVMInstanceRequiredIDs{
			InstanceIds: &[]string{*vm.InstanceID},
		})
		if err != nil {
			return &api.PluginError{Err: err, Step: api.PluginStepUpdateMasterAgentPoolUpdateVMs}
		}

		u.log.Infof("reimaging %s", hostname)
		err = u.vmc.Reimage(ctx, cs.Properties.AzProfile.ResourceGroup, ssName, *vm.InstanceID, nil)
		if err != nil {
			return &api.PluginError{Err: err, Step: api.PluginStepUpdateMasterAgentPoolReimage}
		}

		u.log.Infof("starting %s", hostname)
		err := u.vmc.Start(ctx, cs.Properties.AzProfile.ResourceGroup, ssName, *vm.InstanceID)
		if err != nil {
			return &api.PluginError{Err: err, Step: api.PluginStepUpdateMasterAgentPoolStart}
		}

		u.log.Infof("waiting for %s to be ready", hostname)
		err = u.kubeclient.WaitForReadyMaster(ctx, hostname)
		if err != nil {
			return &api.PluginError{Err: err, Step: api.PluginStepUpdateMasterAgentPoolWaitForReady}
		}

		blob.HostnameHashes[*vm.Name] = desiredHash
		if err := u.updateBlobService.Write(blob); err != nil {
			return &api.PluginError{Err: err, Step: api.PluginStepUpdateMasterAgentPoolUpdateBlob}
		}
	}

	return nil
}
