package v3rpc

import (
	"context"
	"time"

	"go.etcd.io/etcd/api/v3/membershippb"
	pb "go.etcd.io/etcd/api/v3/etcdserverpb"
	"go.etcd.io/etcd/server/v3/etcdserver"
	"go.etcd.io/etcd/api/v3/version"
)

type MaintenanceServer struct {
	rg  etcdserver.RaftKV
	r   etcdserver.RaftStatusGetter
	hdr header
}

func NewMaintenanceServer(s *etcdserver.EtcdServer) pb.MaintenanceServer {
	return &MaintenanceServer{
		rg:  s,
		r:   s,
		hdr: newHeader(s),
	}
}

func (as *MaintenanceServer) Status(ctx context.Context, r *pb.StatusRequest) (*pb.StatusResponse, error) {
	leader := uint64(as.r.Leader())
	if leader != 0 && !as.r.HasLeader() {
		leader = 0
	}
	resp := &pb.StatusResponse{
		Header:           as.hdr.newHeader(),
		Version:          version.Version,
		DbSize:           as.rg.Backend().Size(),
		Leader:           leader,
		RaftIndex:        as.r.CommittedIndex(),
		RaftTerm:         as.r.Term(),
		RaftAppliedIndex: as.r.AppliedIndex(),
	}
	return resp, nil
}

func (as *MaintenanceServer) Alarm(ctx context.Context, r *pb.AlarmRequest) (*pb.AlarmResponse, error) {
	return nil, nil
}

func (as *MaintenanceServer) Defragment(ctx context.Context, r *pb.DefragmentRequest) (*pb.DefragmentResponse, error) {
	return nil, nil
}

func (as *MaintenanceServer) Snapshot(r *pb.SnapshotRequest, stream pb.Maintenance_SnapshotServer) error {
	return nil
}

func (as *MaintenanceServer) MoveLeader(ctx context.Context, r *pb.MoveLeaderRequest) (*pb.MoveLeaderResponse, error) {
	return nil, nil
}
