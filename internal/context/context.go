package context

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/hekmon/transmissionrpc/v2"
	"github.com/notjedi/gotem/internal/config"
)

type Context struct {
	Client    *transmissionrpc.Client
	ListWidth float32
}

var contextInstance *Context

func GetContext(c config.Config) (*Context, error) {
	if contextInstance == nil {
		client, err := transmissionrpc.New(c.Host, c.Username, c.Password,
			&transmissionrpc.AdvancedConfig{
				Port:        c.Port,
				Debug:       c.Debug,
				RPCURI:      c.RpcPath,
				HTTPTimeout: 10 * time.Second,
			})
		if err != nil {
			return nil, err
		}

		ctx, cancel := context.WithCancel(context.TODO())
		defer cancel()

		if ok, serverVersion, serverMinimumVersion, err := client.RPCVersion(ctx); err != nil {
			return nil, err
		} else if !ok {
			return nil, errors.New(fmt.Sprintf(`Remote transmission RPC version (v%d) is
                incompatible with the transmission library (v%d): remote needs at least v%d`,
				serverVersion, transmissionrpc.RPCVersion, serverMinimumVersion))
		}
		contextInstance = &Context{}
		contextInstance.Client = client
	}
	return contextInstance, nil
}
