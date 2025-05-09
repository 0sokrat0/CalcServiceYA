package grpc

import (
	"fmt"

	"github.com/0sokrat0/GoApiYA/agent/config"
	gen "github.com/0sokrat0/GoApiYA/agent/pkg/gen/api/task"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewCalcClient(cfg *config.Config) (gen.TaskServiceClient, error) {
	conn, err := grpc.NewClient(
		fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}

	return gen.NewTaskServiceClient(conn), nil
}
