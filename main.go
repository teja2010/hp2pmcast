package main

import (
	"os"
	"io"
	"fmt"
	"log"
	"net"
	"time"
	"context"
	"crypto/sha1"
	"encoding/json"

	"google.golang.org/grpc"
)

const (
	port = ":60001"

	topFTsize = 64
	midFTsize = 64
	botFTsize = 64
)

var (
	ucsdNode = "" // known ucsd lab machine
	m *mcaster
	pktChan = make(chan FwdPacket, 100)
	ctrlChan = make(chan ctrlMsg, 10)
)

type Configuration struct {
	Hostname string
	MulticastFlag int
	Port string
	XBytes int
	YSeconds int
	//Add new items here. read from config.json
	Threshold []int64
	UcsdNode string
}

type ctrlMsg struct {
	msgType int
	// =1 : run do_getFE(), node at hierarchy h, index i is invalid
	// =2 : update nodeId at h, i with fe
	// =3 : JoinCluster at hierarchy h is complete
	h int
	i int
	fe FingerEntry
}

type mcaster struct {
	UnimplementedMcasterServer

	s *grpc.Server
	self NodeId
	hostname string
	topFT [topFTsize]FingerEntry
	midFT [midFTsize]FingerEntry
	botFT [botFTsize]FingerEntry
	config Configuration
}

type joinListStruct struct {
	heir int
	hostname string
	resp JoinResp
}

func McastHash(s string) uint64 {
	h := sha1.New()
	io.WriteString(h, s)
	b := h.Sum(nil)
	var x uint64
	x = ((uint64(b[0]) << 56) | (uint64(b[1]) << 48) |
	     (uint64(b[2]) << 40) | (uint64(b[3]) << 32) |
	     (uint64(b[4]) << 24) | (uint64(b[5]) << 16) |
	     (uint64(b[6]) << 8) | (uint64(b[7])))
	return x
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

func DoJoin(hostname string, joinChan chan JoinResp) {
	//
}

func ListFindAndPop(ll []string, hs string) []string {
	for index,hh := range ll {
		if hh == hs {
			ll2 := make([]string, 0)
			if index != 0 {
				ll2 = append(ll2, ll[0:index-1]...)
			}
			if index != len(ll)-1 {
				ll2 = append(ll2, ll[index+1:]...)
			}
			return ll2
		}
	}

	return ll
}

func NodeIdLessThanEq(id1 , id2 NodeId) bool {
	//TODO: not simple less than ??
	if id1.A < id2.A {
		return true
	}
	if id1.B < id2.B {
		return true
	}
	if id1.C < id2.C {
		return true
	}
	if id1.A == id2.A && id1.B == id2.B && id1.C == id2.C {
		return true
	}

	return false
}

func HalfId(n NodeId) NodeId {
	var half NodeId
	topBit := uint64(1) << 63;
	half.A = n.A ^ topBit
	half.B = n.B ^ topBit
	half.C = n.C ^ topBit
	return half
}

func NodeIdBetween(id1, left, right NodeId) bool {
	// is left < id1 < right ?

	return false
}

func JoinCluster(m *mcaster, hierarchy int, rootNode string) {
	if hierarchy > 2 {
		return
	}
	h := McastHash(rootNode)
	rootId := NodeId{A:h, B:h, C:h}
	rootFE := FingerEntry{Id:&rootId, Hostname:rootNode}
	if m.hostname == ucsdNode { // also == rootNode
		// we are the first one. set everything to rootId
		for i,_ := range m.topFT {
			m.topFT[i] = rootFE
		}
		for i,_ := range m.midFT {
			m.midFT[i] = rootFE
		}
		for i,_ := range m.botFT {
			m.botFT[i] = rootFE
		}
		return
	}
	halfId := HalfId(rootId)
	// code to join cluster
	joinChan := make(chan JoinResp, 100)
	lowerClusterFound := false

	// nodes between Root & Half.
	// nodes between Half & Root. Clockwise direction
	openReqsR2H := make([]string, 0)
	openReqsH2R := make([]string, 0)

	openReqsR2H = append(openReqsR2H, rootNode)
	go DoJoin(rootNode, joinChan)

	for len(openReqsR2H) > 0 {
	select {
	case jr := <-joinChan:
		//send the join req, add it maybe?
		ctrlChan <- ctrlMsg{msgType:1, h:hierarchy, i:-1, fe:*jr.Self}
		openReqsR2H = ListFindAndPop(openReqsR2H, jr.Self.Hostname)

		//check threshold start lower join cluster request if possible
		rtt_avg := (jr.RttMs + (Time64() - jr.Time)/1000000)
		if !lowerClusterFound && rtt_avg < m.config.Threshold[h] {
			go JoinCluster(m, hierarchy +1, jr.Self.Hostname)
			lowerClusterFound = true
		}

		for _, fe := range jr.FEList {
			if !NodeIdBetween(*fe.Id, rootId, halfId) {
				openReqsH2R = append(openReqsH2R, fe.Hostname)
				continue
			}

			openReqsR2H = append(openReqsR2H, fe.Hostname)
			go DoJoin(fe.Hostname, joinChan)
		}
		//case <-time.After(1 * time.Second):
		//	//check len before waiting
	}
	}

	// same thing for openReqsH2R. but dont add in openReqR2H
	for len(openReqsH2R) > 0 {
	select {
	case jr := <-joinChan:
		//send the join req, add it maybe?
		ctrlChan<-ctrlMsg{msgType:1, h:hierarchy, i:-1, fe:*jr.Self}
		openReqsH2R = ListFindAndPop(openReqsH2R, jr.Self.Hostname)

		//check threshold start lower join cluster request if possible
		rtt_avg := (jr.RttMs + (Time64() - jr.Time)/1000000)
		if !lowerClusterFound && rtt_avg < m.config.Threshold[h] {
			go JoinCluster(m, hierarchy + 1, jr.Self.Hostname)
			lowerClusterFound = true
		}

		for _, fe := range jr.FEList {
			if !NodeIdBetween(*fe.Id, rootId, halfId) {
				continue
			}

			openReqsH2R = append(openReqsH2R, fe.Hostname)
			go DoJoin(fe.Hostname, joinChan)
		}
	}
	}

	if !lowerClusterFound {
		// start own cluster. TODO
	}

	ctrlChan<-ctrlMsg{msgType:3, h:hierarchy} //signal completion of join
}

func FillFingerTable(m *mcaster) {
	// code to fill up the finger table
	for {
	select {
	}
	}
}

func DoFwdPkt(fe FingerEntry, pkt FwdPacket, hierarchy int) {
	//call rpc to fwd it
	log.Printf("Enter DoFwdPkt to %v.%v.%v", fe.Id.A, fe.Id.B, fe.Id.C)

	conn, err := grpc.Dial(fe.Hostname, grpc.WithInsecure(), grpc.WithBlock());
	if err!= nil {
		log.Printf("Dial failed: fe %v", fe.Hostname)
		return
	}
	defer conn.Close()

	new_pkt := new(FwdPacket)
	hops := pkt.EvalList[len(pkt.EvalList)-1].Hops + 1
	new_pkt.EvalList = append(new_pkt.EvalList, pkt.EvalList...)
	new_pkt.EvalList = append(pkt.EvalList, &EvalInfo{Hops: hops, Time: Time64()})
	new_pkt.Payload = append(new_pkt.Payload, pkt.Payload...)

	c := NewMcasterClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 500 * time.Millisecond) //wait for 500ms
	defer cancel()
	_, err = c.Fwd(ctx, new_pkt)
	if err!= nil {
		log.Printf("Fwd to %v failed", fe.Hostname)
		//send a message to ctrlChan
		return
	}
	log.Printf("FWD to %v complete", fe.Hostname)

}

func Time64() int64 {
	// returns nanoseconds since 1970
	// see https://golang.org/pkg/time/#Time.UnixNano
	ts := int64(time.Now().UnixNano())
	return ts
}

func NeedToFwd(fe FingerEntry, pkt FwdPacket, h int) bool {
	s := *pkt.Src
	//l := *pkt.Limit
	if fe.Id == nil { //unfilled fingerEntry
		return false
	}

	
	//ERROR: invalid indirect of pkt.Src.B (type int64) ?? WHATISTHIS
	if s.A == m.self.A && s.B == m.self.B && s.C == m.self.C {
		// TODO check if this is a repeated entry. dont send
		return true // we are are the source. fwd to everyone we know
	}

	if s.A != m.self.A {
		// diff top level hierarchy:
		//      - fwd in top level within limit
		//      - fwd to all in bottom and lower levels
		if h == 1 && fe.Id.A < pkt.Limit.A {
			return true
		}
		if h > 1 {
			return true
		}

		return false
	} else if pkt.Src.B != m.self.B {
		// same top level, but diff mid level hierarchy:
		//      - fwd in mid level within limit
		//      - fwd to all in bottom and lower levels
		if h == 2 && fe.Id.B < pkt.Limit.B {
			return true
		}
		if h > 2 {
			return true
		}

		return false
	} else {
		// same top and mid level hierarchy:
		//  - fwd in bottom within limit
		if h == 3 && fe.Id.C < pkt.Limit.C {
			return true
		}

		return false
	}

}

func ReadConfig(file string) Configuration {
	fp, err := os.Open(file)
	if err != nil {
		log.Fatalf("Config file open failed %v", err)
	}
	defer fp.Close()

	decoder := json.NewDecoder(fp)
	c := Configuration{}
	err = decoder.Decode(&c)
	if err != nil {
		log.Fatalf("Config file read failed: %v", err)
	}

	log.Printf("Configuration: MulticastFlag %v, Port %v",
			c.MulticastFlag, c.Port)
	return c
}

func StartMcast() {
	for {
		log.Printf("Mcast start")
		//send out X bytes
		xbytes := make([]byte, m.config.XBytes)//zeros
		pkt := new(FwdPacket)
		pkt.Payload = append(pkt.Payload, xbytes...)
		self := m.self
		pkt.Src = &self
		pktChan <- *pkt

		// sleep for Y seconds
		time.Sleep(time.Duration(m.config.YSeconds) * time.Second)
		log.Printf("Mcast Complete")
	}
}

func main() {
	fmt.Println("--- Hierarchical P2P Mcast ---")

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalln("failed to open tcp socket: %v", err)
	}

	m = new(mcaster)
	m.s = grpc.NewServer()
	m.config = ReadConfig(string("config.json"))
	m.hostname = m.config.Hostname
	ucsdNode = m.config.UcsdNode
	RegisterMcasterServer(m.s, m)

	go ServeRPCServer(m.s, lis)

	go JoinCluster(m, 1, ucsdNode)
	FillFingerTable(m)

	if m.config.MulticastFlag == 1 {
		go StartMcast()
	}

	for {
	select {

	case pkt := <-pktChan:
		log.Printf("received pkt")
		//decide fingertable based on m.self and pkt.src
		for _, fe := range m.topFT {
			if NeedToFwd(fe, pkt, 1) {
				go DoFwdPkt(fe, pkt, 1)
			}
		}
		for _, fe := range m.midFT {
			if NeedToFwd(fe, pkt, 2) {
				go DoFwdPkt(fe, pkt, 2)
			}
		}
		for _, fe := range m.botFT {
			if NeedToFwd(fe, pkt, 3) {
				go DoFwdPkt(fe, pkt, 3)
			}
		}

	// to handle failures
	//case cmsg := <-ctrlChan:
		// TODO: handle failures
		//if cmsg.msgType == 1 {
		//	if cmsg.h == 1 {
		//		go DoGetFE(m.topFT[i], )
		//	}
		//	go DoGetFE(cmsg.h, cmsg.i)
		//} else if cmsg.msgType == 2 {
		//	if cmsg.h == 1 {
		//		m.topFT[i] = fe
		//	}
		//}
	}
	}
}
