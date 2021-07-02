package vsphere

import (
	"context"
	"log"
	"net/url"

	"github.com/vmware/govmomi/session/cache"
	"github.com/vmware/govmomi/vim25"
	"github.com/vmware/govmomi/vim25/soap"
)

type VshpereClient struct {
	C   *vim25.Client
	S   *cache.Session
	Ctx context.Context
}

var NewVshpereClient = &VshpereClient{}

func init() {
	ctx := context.Background()
	u, err := soap.ParseURL("192.168.4.240") // https://10.50.82.155/sdk
	u.User = url.UserPassword("administrator@51elab.local", "Hello51elab.manager")

	// Share govc's session cache
	s := &cache.Session{
		URL:      u,
		Insecure: true,
	}

	c := new(vim25.Client)
	err = s.Login(ctx, c, nil)
	if err != nil {
		log.Fatal(err)
	}

	// defer s.Logout(ctx, c)

	if err != nil {
		log.Fatal(err)
	}

	NewVshpereClient.C = c
	NewVshpereClient.S = s
	NewVshpereClient.Ctx = ctx
}
