package context

import (
	"context"
	"fmt"
	"time"

	"github.com/hekmon/transmissionrpc/v2"
	"github.com/notjedi/gotem/internal/config"
)

type Context interface {
	Client() *transmissionrpc.Client
	TitleSpacing() [3]uint
	DescSpacing() [4]uint
	SetTitleSpacing([3]uint)
	SetDescSpacing([4]uint)
}

type ProgramContext struct {
	client       *transmissionrpc.Client
	titleSpacing [3]uint
	descSpacing  [4]uint
}

func (c *ProgramContext) Client() *transmissionrpc.Client {
	return c.client
}

func (c *ProgramContext) TitleSpacing() [3]uint {
	return c.titleSpacing
}

func (c *ProgramContext) SetTitleSpacing(titleSpacing [3]uint) {
	// TODO: should i be using a lock here?
	c.titleSpacing = titleSpacing
}

func (c *ProgramContext) DescSpacing() [4]uint {
	return c.descSpacing
}

func (c *ProgramContext) SetDescSpacing(descSpacing [4]uint) {
	// TODO: should i be using a lock here?
	c.descSpacing = descSpacing
}

var contextInstance Context

// TODO: sync this across threads - so there is only one Context object globally
func GetContext(c config.Config) (Context, error) {
	if contextInstance == nil {
		client, err := transmissionrpc.New(c.Host, c.Username, c.Password,
			&transmissionrpc.AdvancedConfig{
				Port:        c.Port,
				Debug:       c.Debug,
				RPCURI:      c.RPCPath,
				HTTPTimeout: 10 * time.Second,
			})
		if err != nil {
			return nil, err
		}

		// TODO: why is this even here? should i really need a cancel callback here?
		ctx, cancel := context.WithCancel(context.TODO())
		defer cancel()

		if ok, serverVersion, serverMinimumVersion, err := client.RPCVersion(ctx); err != nil {
			return nil, err
		} else if !ok {
			return nil, fmt.Errorf(`Remote transmission RPC version (v%d) is
                incompatible with the transmission library (v%d): remote needs at least v%d`,
				serverVersion, transmissionrpc.RPCVersion, serverMinimumVersion)
		}
		// TODO: no need for pointer to interface? anyways pointer to interface seems odd?
		// https://github.com/uber-go/guide/blob/master/style.md#pointers-to-interfaces
		// https://stackoverflow.com/questions/54670125/how-to-get-pointer-to-interface-in-go/54670253#54670253
		contextInstance = Context(&ProgramContext{client: client})
	}
	return contextInstance, nil
}
