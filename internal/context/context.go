package context

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/hekmon/transmissionrpc/v2"
	c "github.com/notjedi/gotem/internal/config"
)

var client *transmissionrpc.Client

func GetClient(config c.Config) (*transmissionrpc.Client, error) {
	if client == nil {
		client, err := transmissionrpc.New(config.Host, config.Username, config.Password,
			&transmissionrpc.AdvancedConfig{
				Port:        config.Port,
				Debug:       config.Debug,
				RPCURI:      config.RpcPath,
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
			client = nil
			return nil, errors.New(fmt.Sprintf(`Remote transmission RPC version (v%d) is
                incompatible with the transmission library (v%d): remote needs at least v%d`,
				serverVersion, transmissionrpc.RPCVersion, serverMinimumVersion))
		}
	}
	return client, nil
}
