package clientv3

import (
	"context"
	"time"

	pb "go.etcd.io/etcd/api/v3/etcdserverpb"
	"google.golang.org/grpc"
)

type Maintenance interface {
	Status(ctx context.Context, endpoint string) (*StatusResponse, error)
}

type maintenance struct {
	c *Client
}

type StatusResponse pb.StatusResponse

func NewMaintenance(c *Client) Maintenance {
	return &maintenance{c: c}
}

func (m *maintenance) Status(ctx context.Context, endpoint string) (*StatusResponse, error) {
	conn, err := grpc.DialContext(ctx, endpoint, m.c.cfg.DialOptions...)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	cli := pb.NewMaintenanceClient(conn)
	resp, err := cli.Status(ctx, &pb.StatusRequest{})
	if err != nil {
		return nil, err
	}
	return (*StatusResponse)(resp), nil
}
