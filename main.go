package main

import (
	"fmt"
	"net"
	"time"
	"context"

	"google.golang.org/grpc"
)

const (
	port = ":60001"

	// 128 = 32 + 32 + 64
	topFTsize = 32
	midFTsize = 32
	botFTsize = 64
	knownNode = "google.com" // lol. some known node
)


type mcaster struct {
	UnimplementedMcasterServer

	s *grpc.Server
	topFT [topFTsize]FingerEntry
	midFT [midFTsize]FingerEntry
	botFT [botFTsize]FingerEntry
}

func (s *mcaster) Fwd(ctx context.Context, in *FwdPacket) (*Empty, error) {
	fmt.Println("Fwd Limit %v.%v.%v", in.Limit.A, in.Limit.B, in.Limit.C)
	out := new(Empty)

	return out, nil
}

func (s *mcaster) Join(ctx context.Context, in *JoinReq) (*JoinResp, error) {
	fmt.Println("JoinReq Hierarchy %v", in.Hierarchy)
	out := new(JoinResp)
	out.Hierarchy = in.Hierarchy

	return out, nil
}

func (s *mcaster) GetFingerEntry(ctx context.Context, in *GetFERequest) (*GetFEResponse, error) {
	fmt.Println("GetFingerEntry Hierarchy %v", in.Hierarchy)
	out := new(GetFEResponse)
	out.Hierarchy = in.Hierarchy

	return out, nil
}

func (s *mcaster) SetSuccessor(ctx context.Context, in *Successor) (*Empty, error) {
	fmt.Println("SetSuccessor Id %v", in.Id)
	out := new(Empty)

	return out, nil
}

func ServeRPCServer (s *grpc.Server, lis net.Listener) {
	err := s.Serve(lis);
	if err != nil {
		fmt.Println("Failed to start server %v", err)
	}
}

func JoinCluster(m *mcaster) {
	// code to join cluster
}

func FillFingerTable(m *mcaster) {
	// code to fill up the finger table
}

func main() {
	fmt.Println("Hierarchical P2P Mcast")

	lis, err := net.Listen("tcp", port)
	if err != nil {
		fmt.Println("failed to open tcp socket: %v", err)
	}

	m := new(mcaster)
	m.s = grpc.NewServer()
	RegisterMcasterServer(m.s, m)

	go ServeRPCServer(m.s, lis)

	JoinCluster(m)
	FillFingerTable(m)

	for {
		// read config file and send out mcast data.
		time.Sleep(1 * time.Second)
		fmt.Println("asd")
	}
}
