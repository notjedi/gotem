package context

import (
	"context"
	"errors"
	"fmt"

	"github.com/hekmon/transmissionrpc/v2"
	c "github.com/notjedi/gotem/internal/config"
)

var client *transmissionrpc.Client

func GetClient(config c.Config) (*transmissionrpc.Client, error) {
	if client == nil {
		client, err := transmissionrpc.New(config.Host, config.Username, config.Password,
			&transmissionrpc.AdvancedConfig{
				Port: config.Port,
			})
		if err != nil {
			return nil, err
		}
		if ok, serverVersion, serverMinimumVersion, err := client.RPCVersion(context.TODO()); err != nil {
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

// func Must(client *transmissionrpc.Client) {
// 	if ok, serverVersion, serverMinimumVersion, err := client.RPCVersion(); err != nil {
// 		return nil, err
// 	} else if !ok {
// 		client = nil
// 		return nil, errors.New(fmt.Sprintf(`Remote transmission RPC version (v%d) is
//                 incompatible with the transmission library (v%d): remote needs at least v%d`,
// 			serverVersion, transmissionrpc.RPCVersion, serverMinimumVersion))
// 	}
// }
