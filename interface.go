package blueprint

import (
	"github.com/hashicorp/go-plugin"
	"net/rpc"
)

type GeneratorRPC struct{ client *rpc.Client }

// Generate makes an RPC client call to Generate
func (c *GeneratorRPC) Generate(spec Spec, values map[string]string, workdir string) error {
	args := &GenerateArgs{
		Spec:    spec,
		Values:  values,
		Workdir: workdir,
	}
	resp := &GenerateResponse{}
	err := c.client.Call("Plugin.Generate", args, resp)
	if err != nil {
		return err
	}
	return resp.Error
}

// LoadSpec makes an RPC client call to LoadSpec
func (c *GeneratorRPC) LoadSpec() (Spec, error) {
	args := &LoadSpecArgs{}
	resp := &LoadSpecResponse{}
	err := c.client.Call("Plugin.LoadSpec", args, resp)
	if err != nil {
		return Spec{}, err
	}
	return resp.Spec, resp.Error
}

// GeneratorRPCServer is the server side of the RPC
type GeneratorRPCServer struct {
	Impl Generator
}

// Generate handles the server-side RPC call for Generate
func (s *GeneratorRPCServer) Generate(args *GenerateArgs, resp *GenerateResponse) error {
	err := s.Impl.Generate(args.Spec, args.Values, args.Workdir)
	resp.Error = err
	return nil
}

// LoadSpec handles the server-side RPC call for LoadSpec
func (s *GeneratorRPCServer) LoadSpec(_ *LoadSpecArgs, resp *LoadSpecResponse) error {
	spec, err := s.Impl.LoadSpec()
	resp.Spec = spec
	resp.Error = err
	return nil
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

type GeneratorPlugin struct {
	Impl Generator
}

func (p *GeneratorPlugin) Server(*plugin.MuxBroker) (interface{}, error) {
	return &GeneratorRPCServer{Impl: p.Impl}, nil
}

func (p *GeneratorPlugin) Client(_ *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &GeneratorRPC{client: c}, nil
}

var HandshakeConfig = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "BASIC_PLUGIN",
	MagicCookieValue: "hello",
}
