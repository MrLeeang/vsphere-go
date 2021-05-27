package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"vsphere-go/pkg/vsphere"

	"github.com/vmware/govmomi/object"

	"github.com/vmware/govmomi/property"

	"github.com/vmware/govmomi/view"

	"github.com/vmware/govmomi/vim25/types"

	"github.com/vmware/govmomi/find"
	"github.com/vmware/govmomi/guest"
	"github.com/vmware/govmomi/guest/toolbox"
	"github.com/vmware/govmomi/vim25/mo"
)

func cloneVm() {
	finder := find.NewFinder(vsphere.NewVshpereClient.C, false)
	defualtDc, err := finder.DefaultDatacenter(vsphere.NewVshpereClient.Ctx)
	if err != nil {
		fmt.Println(err)
		return
	}
	finder.SetDatacenter(defualtDc)

	vm, err := finder.VirtualMachine(vsphere.NewVshpereClient.Ctx, `ubuntu-18.04`)
	if err != nil {
		fmt.Println(err)
		return
	}
	// *object.VirtualMachine 转化为 mo.VirtualMachine
	var vmIns mo.VirtualMachine

	pc := property.DefaultCollector(vsphere.NewVshpereClient.C)
	// 如果想要全部属性，可以传一个空的字串切片
	err = pc.RetrieveOne(vsphere.NewVshpereClient.Ctx, vm.Reference(), []string{}, &vmIns)

	if err != nil {
		fmt.Println(err)
		return
	}

	// fmt.Println(vmIns.Summary.Config.Name)

	// mo 转换为 *object
	oVmParent := object.NewFolder(vsphere.NewVshpereClient.C, vmIns.Parent.Reference())

	snapshot := vmIns.Snapshot.RootSnapshotList[0]

	cloneSpec := types.VirtualMachineCloneSpec{}
	cloneSpec.Template = false
	cloneSpec.PowerOn = false
	cloneSpec.Snapshot = &snapshot.Snapshot
	cloneSpec.Location.DiskMoveType = "createNewChildDiskBacking"
	// cloneSpec.Config.DeviceChange = []types.BaseVirtualDeviceConfigSpec{}

	_, err = vm.Clone(vsphere.NewVshpereClient.Ctx, oVmParent, "test2", cloneSpec)
}

func runCmd() {
	finder := find.NewFinder(vsphere.NewVshpereClient.C, false)

	// 获取数据中心
	// dc, err := finder.Datacenter(vsphere.NewVshpereClient.Ctx, "/")

	defualtDc, err := finder.DefaultDatacenter(vsphere.NewVshpereClient.Ctx)
	if err != nil {
		fmt.Println(err)
		return
	}
	finder.SetDatacenter(defualtDc)

	vm, err := finder.VirtualMachine(vsphere.NewVshpereClient.Ctx, `ubuntu-18.04`)
	if err != nil {
		fmt.Println(err)
		return
	}

	o := guest.NewOperationsManager(vsphere.NewVshpereClient.C, vm.Reference())
	pm, err := o.ProcessManager(vsphere.NewVshpereClient.Ctx)
	if err != nil {
		fmt.Println(err)
		return
	}
	fm, err := o.FileManager(vsphere.NewVshpereClient.Ctx)
	if err != nil {
		fmt.Println(err)
		return
	}

	family := ""
	var vmIns mo.VirtualMachine
	err = vm.Properties(context.Background(), vm.Reference(), []string{}, &vmIns)
	if err != nil {
		fmt.Println(err)
		return
	}

	if vmIns.Guest != nil {
		family = vmIns.Guest.GuestFamily
	}

	type AuthFlag struct {
		auth types.NamePasswordAuthentication
		proc bool
	}

	var Auth types.BaseGuestAuthentication
	auth := &types.NamePasswordAuthentication{}
	GuestId := strings.ToUpper(vmIns.Config.GuestId)
	if strings.Contains(GuestId, "LINUX") {
		auth.Username = "root"
		auth.Password = "12345678"
	} else {
		auth.Username = "administrator"
		auth.Password = "12345678"
	}
	auth.GuestAuthentication.InteractiveSession = false

	Auth = auth

	c := &toolbox.Client{
		ProcessManager: pm,
		FileManager:    fm,
		Authentication: Auth,
		GuestFamily:    types.VirtualMachineGuestOsFamily(family),
	}

	ecmd := &exec.Cmd{
		// Path: "echo root:12345678 |chpasswd",
		Path: "ifconfig",
		Args: []string{},
		// Env:    cmd.vars,
		// Dir:    cmd.dir,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}

	err = c.Run(vsphere.NewVshpereClient.Ctx, ecmd)
	if err != nil {
		fmt.Println("Run..")
		fmt.Println(err)
		return
	}
}

func uploadVm() {
	// var r io.Reader = os.Stdin

	// f, err := os.Open("1.txt")

	// defer f.Close()

	// r = f

	// cmdAttr := types.GuestPosixFileAttributes{}

	// _ = c.Upload(ctx, r, "/1.txt", soap.DefaultUpload, cmdAttr, true)

	// =====================
}

func getVms() {

	m := view.NewManager(vsphere.NewVshpereClient.C)

	ctx := context.Background()

	v, err := m.CreateContainerView(ctx, vsphere.NewVshpereClient.C.ServiceContent.RootFolder, []string{"VirtualMachine"}, true)
	if err != nil {
		log.Fatal(err)
	}

	defer v.Destroy(ctx)

	// Retrieve summary property for all machines
	// Reference: http://pubs.vmware.com/vsphere-60/topic/com.vmware.wssdk.apiref.doc/vim.VirtualMachine.html
	var vms []mo.VirtualMachine
	err = v.Retrieve(ctx, []string{"VirtualMachine"}, []string{}, &vms)
	if err != nil {
		log.Fatal(err)
	}

	// Print summary per vm (see also: govc/vm/info.go)

	for _, vm := range vms {
		// fmt.Printf("%s: %s\n", vm.Summary.Config.Name, vm.Summary.Config.GuestFullName)
		fmt.Printf("%s\n", vm.Summary.Config.Name)
		if vm.Guest != nil {
			fmt.Println(vm.Guest.GuestState)
		}

		// mo.VirtualMachine 转换为 *object.VirtualMachine
		// ovm := object.NewVirtualMachine(vsphere.NewVshpereClient.C, vm.Reference())
		// ovm.PowerOff(vsphere.NewVshpereClient.Ctx)

	}

}

func main() {

	defer vsphere.NewVshpereClient.S.Logout(vsphere.NewVshpereClient.Ctx, vsphere.NewVshpereClient.C)
	runCmd()
	// cloneVm()
	// getVms()
}
