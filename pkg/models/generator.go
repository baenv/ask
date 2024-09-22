package models

import (
	"fmt"
	"sync"

	"github.com/bwmarrin/snowflake"
)

const (
	userNodeID              = 1
	userAgentConfigNodeID   = 2
	serverNodeID            = 3
	serverAdminConfigNodeID = 4
)

var (
	userIDGenerator              *snowflake.Node
	userAgentConfigIDGenerator   *snowflake.Node
	serverIDGenerator            *snowflake.Node
	serverAdminConfigIDGenerator *snowflake.Node
	once                         sync.Once
)

// InitIDGenerators initializes all snowflake node generators
func InitIDGenerators() error {
	var err error
	once.Do(func() {
		userIDGenerator, err = snowflake.NewNode(userNodeID)
		if err != nil {
			err = fmt.Errorf("failed to initialize user ID generator: %w", err)
			return
		}

		userAgentConfigIDGenerator, err = snowflake.NewNode(userAgentConfigNodeID)
		if err != nil {
			err = fmt.Errorf("failed to initialize user agent config ID generator: %w", err)
			return
		}

		serverIDGenerator, err = snowflake.NewNode(serverNodeID)
		if err != nil {
			err = fmt.Errorf("failed to initialize server ID generator: %w", err)
			return
		}

		serverAdminConfigIDGenerator, err = snowflake.NewNode(serverAdminConfigNodeID)
		if err != nil {
			err = fmt.Errorf("failed to initialize server admin config ID generator: %w", err)
			return
		}
	})
	return err
}
