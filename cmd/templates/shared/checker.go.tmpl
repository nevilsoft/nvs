package shared

import (
	"net/rpc"

	"github.com/hashicorp/go-plugin"
)

// Checker interface ที่เดิม
type Checker interface {
	Check() bool
}

// CheckerRPC - RPC wrapper สำหรับ client
type CheckerRPC struct{ client *rpc.Client }

func (c *CheckerRPC) Check() bool {
	var result bool
	err := c.client.Call("Plugin.Check", new(interface{}), &result)
	if err != nil {
		return false
	}
	return result
}

// CheckerRPCServer - RPC server wrapper
type CheckerRPCServer struct {
	Impl Checker
}

func (s *CheckerRPCServer) Check(args interface{}, resp *bool) error {
	*resp = s.Impl.Check()
	return nil
}

// CheckerPlugin - plugin implementation
type CheckerPlugin struct {
	Impl Checker
}

func (p *CheckerPlugin) Server(*plugin.MuxBroker) (interface{}, error) {
	return &CheckerRPCServer{Impl: p.Impl}, nil
}

func (CheckerPlugin) Client(b *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &CheckerRPC{client: c}, nil
}
