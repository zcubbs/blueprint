package blueprint

import (
	"github.com/hashicorp/go-plugin"
	"net/rpc"
)

// RPCPlugin wraps a blueprint.Generator for RPC communication
type RPCPlugin struct {
	Impl Generator
}

func (b *RPCPlugin) Generate(args *GenerateArgs, resp *GenerateResponse) error {
	err := b.Impl.Generate(args.Spec, args.Values, args.Workdir)
	resp.Error = err
	return nil
}

func (b *RPCPlugin) LoadSpec(_ *LoadSpecArgs, resp *LoadSpecResponse) error {
	spec, err := b.Impl.LoadSpec()
	resp.Spec = spec
	resp.Error = err
	return nil
}

// Server returns an RPC server
func (b *RPCPlugin) Server(*plugin.MuxBroker) (interface{}, error) {
	return &RPCServer{Impl: b.Impl}, nil
}

// Client returns an RPC client
func (b *RPCPlugin) Client(_ *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &RPCClient{client: c}, nil
}

// RPCServer is the server side of the RPC
type RPCServer struct {
	Impl Generator
}

// RPCClient is the client side of the RPC
type RPCClient struct {
	client *rpc.Client
}

// RPCPluginClient handles client-side communication for blueprint plugins
type RPCPluginClient struct {
	client *plugin.Client
}

// GenerateArgs Define argument and response structures for the methods
type GenerateArgs struct {
	Spec    Spec
	Values  map[string]string
	Workdir string
}

type GenerateResponse struct {
	Error error
}

type LoadSpecArgs struct {
	// Arguments (if any)
}

type LoadSpecResponse struct {
	Spec  Spec
	Error error
}
