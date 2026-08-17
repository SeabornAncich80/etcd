package clientv3_test

import (
	"context"
	"testing"
	"time"

	"go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/tests/v3/integration"
)

func TestStatusTransientLeader(t *testing.T) {
	integration.BeforeTest(t)
	clus := integration.NewClusterV3(t, &integration.ClusterConfig{Size: 3})
	defer clus.Terminate(t)

	cli := clus.RandClient()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Get initial leader
	resp, err := cli.Status(ctx, clus.Members[0].GRPCAddr())
	if err != nil {
		t.Fatal(err)
	}
	origLeader := resp.Leader

	// Stop the leader node
	var leaderIdx int
	for i, m := range clus.Members {
		resp, err := cli.Status(ctx, m.GRPCAddr())
		if err == nil && resp.Header.MemberId == origLeader {
			leaderIdx = i
			break
		}
	}
	clus.Members[leaderIdx].Stop(t)

	// Wait for election
	time.Sleep(3 * time.Second)

	// Verify status on remaining nodes
	for i, m := range clus.Members {
		if i == leaderIdx {
			continue
		}
		resp, err := cli.Status(ctx, m.GRPCAddr())
		if err != nil {
			t.Fatal(err)
		}
		if resp.Leader == origLeader {
			t.Errorf("node %d still reports stale leader %d", i, origLeader)
		}
		if resp.Leader == 0 {
			t.Errorf("node %d reports no leader", i)
		}
	}
}
