package grpchandlers

import (
	context "context"

	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

func (s *GRPCServer) Ping(ctx context.Context, in *PingRequest) (*PingResponse, error) {
	err := s.storage.Ping()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Database is not working")
	}

	return &PingResponse{Result: "OK"}, nil
}
