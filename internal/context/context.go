package context

import (
	"context"
	"fmt"
	"time"

	"github.com/hekmon/transmissionrpc/v2"
	"github.com/notjedi/gotem/internal/config"
)

type Context struct {
	Client       *transmissionrpc.Client
	ListWidth    float32
	TitleSpacing [3]uint
	DescSpacing  [4]uint
}

var contextInstance *Context

// TODO: sync this across threads - so there is only one Context object globally
func GetContext(c config.Config) (*Context, error) {
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
		contextInstance = &Context{}
		contextInstance.Client = client
	}
	return contextInstance, nil
}
