package snowflake

import (
	"context"

	"github.com/bwmarrin/snowflake"
	"github.com/gosharedlib/idgenerator/workid"
)

const defaultEpoch = 1648656000000

type snowflakeIDGenerator struct {
	node *snowflake.Node
}

// Generator ID生成器
type Generator interface {
	// GenID 生成字符串 Key.
	GenID() string
	// GenIntID 生成整型 Key.
	GenIntID() int64
}

func NewSnowflakeGenerator(worker workid.Conn, epoch ...int64) Generator {
	snowflake.Epoch = defaultEpoch
	if len(epoch) > 0 {
		snowflake.Epoch = epoch[0]
	}

	workID, err := worker.GetWorkID(context.Background())
	if err != nil {
		panic(err)
	}
	node, err := snowflake.NewNode(int64(workID))
	if err != nil {
		panic(err)
	}
	if err != nil {
		panic(err)
	}

	return &snowflakeIDGenerator{node: node}
}

func (g snowflakeIDGenerator) GenID() string {
	return g.node.Generate().String()
}

func (g snowflakeIDGenerator) GenIntID() int64 {
	return g.node.Generate().Int64()
}
