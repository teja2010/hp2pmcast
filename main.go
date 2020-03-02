package main

import (
	"os"
	"io"
	"fmt"
	"log"
	"net"
	"time"
	"errors"
	"context"
	"runtime"
	"strconv"
	"math/rand"
	"crypto/sha1"
	"encoding/json"

	"google.golang.org/grpc"
)

const (
	NANO_TO_MILLISEC = 1000000   // 10^6
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
	XBytes int
	YSeconds int
	//Add new items here. read from config.json
	Threshold []int64
	UcsdNode string
	Dlog int
	LogColor int
}
func (c Configuration) String() string {
	return fmt.Sprintf("{Hostname %v, MulticastFlag %v, " +
			   "XBytes %v, YSeconds %v, Threshold %v, UcsdNode %v}",
			   c.Hostname, c.MulticastFlag,
			   c.XBytes, c.YSeconds,
			   c.Threshold, c.UcsdNode)
}
func GetHostName(hostname string) string {
	for i:=len(hostname)-1 ; i> 0 ;i-- {
		if hostname[i] == ':' {
			return hostname[i:]
		}
	}
	return hostname
}

const (
	MSG_TYPE_RUN_GETFE       = iota  //0
	MSG_TYPE_ADD_NODEID
	MSG_TYPE_JOIN_COMPLETE
	MSG_TYPE_FAILED_FE
)

// once a node hostname and Id are found (only in top and mid lvl)
// 1. send a GetFingerEntry to get an alternate host in that cluster.
//    The receiver will trigger a GetFingerEntry if the cluster is unknown
//    and could be added to the finger table.
// 2. on recving a response update, add to finger table.
// use MSG_TYPE_RUN_GETFE and MSG_TYPE_ADD_NODEID to do this

// in low lvl, once a node hostname is found,
// 1. send a set succesor/set predecessor. similar to chord

type ctrlMsg struct {
	msgType int
	// =1 : run do_getFE(), node at hierarchy, index i is invalid
	// =2 : update nodeId at hierarchy, i with fe
	// =3 : JoinCluster at hierarchy is complete
	// =4 : Sending to fe has failed, replace the finger entry
	hierarchy int
	fe FingerEntry
}
func (c ctrlMsg) String() string {
	return fmt.Sprintf("{ msgType %v, hierarchy %v, " +
			   "FingerEntry %v}",
			   c.msgType, c.hierarchy, c.fe)
}

type mcaster struct {
	UnimplementedMcasterServer

	s *grpc.Server
	id NodeId
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

func LogFormat() string {
	_, file, line, ok := runtime.Caller(2)
	if !ok {
		file = "???"
		line = 0
	}
	for i := len(file) - 1; i > 0; i-- {
		if file[i] == '/' {
			file = file[i+1:]
			break
		}
	}
	return (file + ":" + strconv.Itoa(line) + " ")
}

func DLog(format string, v ...interface{}) {
	if m.config.Dlog == 0 {
		return
	}

	file := LogFormat()
	if m.config.LogColor == 1{
		log.Printf("\x1b[36m DEBUG | " + file + format + "\033[0m", v ...)
	} else {
		log.Printf(" DEBUG | " +file+ format , v ...)
	}
}

func MLog(format string, v ...interface{}) {
	file := LogFormat()
	log.Printf(" MID   | " + file + format , v ...)
}

func HLog(format string, v ...interface{}) {
	file := LogFormat()
	if m.config.LogColor == 1{
		log.Printf("\x1B[31m HIGH  | " + file+ format + "\033[0m", v ...)
	} else {
		log.Printf(" HIGH  | " +file+ format , v ...)
	}
}

func McastHash(s string) uint64 {
	sha := sha1.New()
	io.WriteString(sha, s)
	b := sha.Sum(nil)
	var x uint64
	x = ((uint64(b[0]) << 56) | (uint64(b[1]) << 48) |
	     (uint64(b[2]) << 40) | (uint64(b[3]) << 32) |
	     (uint64(b[4]) << 24) | (uint64(b[5]) << 16) |
	     (uint64(b[6]) << 8) | (uint64(b[7])))
	return x
}

func (s *mcaster) Fwd(ctx context.Context, in *FwdPacket) (*Empty, error) {
	DLog("Fwd in FwdPacket %v", in)
	out := new(Empty)

	return out, nil
}

func (s *mcaster) Join(ctx context.Context, in *JoinReq) (*JoinResp, error) {
	timeNow := Time64()
	DLog("Join : JoinReq %v", in)
	out := new(JoinResp)
	out.Hierarchy = in.Hierarchy
	out.RttMs = (timeNow - in.Time)/NANO_TO_MILLISEC
	out.Time = timeNow
	out.Self = &FingerEntry{Id:&m.id, Hostname:m.hostname}
	if in.Hierarchy == 1 {
		out.FEList = ValidFE_ptr(m.topFT[:])
	} else if in.Hierarchy == 2 {
		out.FEList = ValidFE_ptr(m.midFT[:])
	} else if in.Hierarchy == 3 {
		out.FEList = ValidFE_ptr(m.botFT[:])
	}
	DLog("JoinResp %v", out)

	return out, nil
}

func (s *mcaster) GetFingerEntry(ctx context.Context, in *GetFERequest) (*GetFEResponse, error) {
	DLog("GetFingerEntry : GetFERequest %v", in)
	out := new(GetFEResponse)
	out.Hierarchy = in.Hierarchy

	return out, nil
}

func (s *mcaster) SetSuccessor(ctx context.Context, in *Successor) (*Empty, error) {
	DLog("SetSuccessor Id %v", in.FE.Id)
	out := new(Empty)

	return out, nil
}

func isEqualNodeId(n1, n2 NodeId) bool {
	return (n1.A == n2.A && n1.B == n2.B && n1.C == n2.C)
}

func isEqualFE(fe1, fe2 FingerEntry) bool {
	return isEqualNodeId(*fe1.Id, *fe2.Id) && fe1.Hostname == fe2.Hostname
}

func ServeRPCServer (s *grpc.Server, lis net.Listener) {
	err := s.Serve(lis);
	if err != nil {
		log.Fatalf("Failed to start server %v", err)
	}
	DLog("Started RPC Server")
}

func DoJoin(hostname string, hierarchy int, joinChan chan JoinResp) {
	MLog("DoJoin hostname %s", hostname)
	conn, err := grpc.Dial(hostname, grpc.WithInsecure(), grpc.WithBlock());
	if err!= nil {
		HLog("Dial failed: fe %v", hostname)
		return
	}
	defer conn.Close()

	c := NewMcasterClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(),
					   500 * time.Millisecond)
					   //wait for 500ms
	defer cancel()
	joinResp, err := c.Join(ctx, &JoinReq{Hierarchy:int32(hierarchy),
					      Time: Time64()})
	if err != nil {
		HLog("Join to %v failed: %v", hostname, err)
		return
	}
	DLog("Join Resp %v", joinResp)
	joinChan<-*joinResp
	DLog("Join Complete")

}

//dont add if hs is already present in ll
func ListFindAndAdd(ll []string, hs string) ([]string, bool) {
	DLog("ListFindAndAdd list: %v, add hostname %v", ll, hs)
	for _,hh := range ll {
		if hh == hs {
			return ll, false
		}
	}

	return append(ll, hs), true
}

func invalidateFE(fe FingerEntry, hierarchy int) {
	//set it self.id and self.hostname. it is invalid

	selfId := new(NodeId)
	*selfId = m.id
	selfFE := FingerEntry{Id:selfId, Hostname:m.hostname}

	switch hierarchy {
	case 1:
		for i := 0; i < topFTsize; i++ {
			if isEqualFE(m.topFT[i], fe) {
				m.topFT[i] = selfFE
			}
		}
	case 2:
		for i := 0; i < midFTsize; i++ {
			if isEqualFE(m.midFT[i], fe) {
				m.midFT[i] = selfFE
			}
		}
	case 3:
		for i := 0; i < botFTsize; i++ {
			if isEqualFE(m.botFT[i], fe) {
				m.botFT[i] = selfFE
			}
		}
	}
}

func getRandomFE(hierarchy int) (FingerEntry, error) {
	selfId := new(NodeId)
	*selfId = m.id
	selfFE := FingerEntry{Id:selfId, Hostname:m.hostname}
	switch hierarchy {
	case 1:
		rr := rand.Intn(topFTsize)
		for i := 0; i < topFTsize; i++ {
			if isEqualFE(m.topFT[(rr + i)%topFTsize], selfFE) {
				return m.topFT[(rr + i)%topFTsize], nil
			}
		}
	case 2:
		rr := rand.Intn(midFTsize)
		for i := 0; i < midFTsize; i++ {
			if isEqualFE(m.midFT[(rr + i)%midFTsize], selfFE) {
				return m.midFT[(rr + i)%midFTsize], nil
			}
		}
	case 3:
		rr := rand.Intn(botFTsize)
		for i := 0; i < botFTsize; i++ {
			if isEqualFE(m.botFT[(rr + i)%botFTsize], selfFE) {
				return m.botFT[(rr + i)%botFTsize], nil
			}
		}
	}

	return selfFE, errors.New("Unable to find a rand FE")
}

func DoGetFE(fe FingerEntry, hierarchy int, limit NodeId) {
	// send a req to fe, @ hierarchy, get an entry within limit

	MLog("Enter DoGetFE fe %v, hierarchy %d, limit %v", fe, hierarchy, limit)

	conn, err := grpc.Dial(fe.Hostname, grpc.WithInsecure(), grpc.WithBlock());
	if err!= nil {
		HLog("Dial failed: fe %v", fe.Hostname)
		return
	}
	defer conn.Close()

	new_fe_req := new(GetFERequest)
	new_fe_req.Hierarchy = int32(hierarchy)
	new_fe_req.Limit = &limit
	DLog("New FE Req %v", new_fe_req)

	c := NewMcasterClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(),
					   500 * time.Millisecond) //wait for 500ms
	defer cancel()
	resp, err := c.GetFingerEntry(ctx, new_fe_req)
	if err != nil {
		HLog("GetFingerEntry to %v failed, err %v", fe.Hostname, err)
		// TODO what to do?
		return
	}
	DLog("GetFingerEntry complete:  %v", resp)

	ctrlChan <- ctrlMsg{msgType:MSG_TYPE_ADD_NODEID,
			    hierarchy : hierarchy,
			    fe : *resp.NewFE,
			   }
}

//find and remove hs from ll
func ListFindAndPop(ll []string, hs string) []string {
	DLog("ListFindAndPop list: %v, pop hostname %v", ll, hs)
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
	DLog("NodeIdBetween: left %v < id1 %v < right %v", left, id1, right)
	// is left < id1 < right ?

	return false
}

func JoinCluster(m *mcaster, hierarchy int, rootNode string) {
	if hierarchy > 2 {
		return
	}
	HLog("JoinCluster hierarchy %d, rootNode %s", hierarchy, rootNode)

	hash := McastHash(rootNode)
	rootId := NodeId{A:hash, B:hash, C:hash}
	if m.hostname == ucsdNode { // also == rootNode
		HLog("The root node. No clusters joined")
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

	openReqsR2H,_ = ListFindAndAdd(openReqsR2H, rootNode)
	go DoJoin(rootNode, hierarchy, joinChan)

	for len(openReqsR2H) > 0 {
	select {
	case jr := <-joinChan:
		//send the join req, add it maybe?
		ctrlChan <- ctrlMsg{msgType:MSG_TYPE_RUN_GETFE,
				    hierarchy:hierarchy,
				    fe:*jr.Self,
				   }
		openReqsR2H = ListFindAndPop(openReqsR2H, jr.Self.Hostname)

		//check threshold start lower join cluster request if possible
		rtt_avg := (jr.RttMs + (Time64() - jr.Time)/NANO_TO_MILLISEC)
		if !lowerClusterFound &&
		   rtt_avg < m.config.Threshold[hierarchy] {
			go JoinCluster(m, hierarchy +1, jr.Self.Hostname)
			lowerClusterFound = true
		}

		for _, fe := range jr.FEList {
			if !NodeIdBetween(*fe.Id, rootId, halfId) {
				openReqsH2R, _ = ListFindAndAdd(openReqsH2R,
								fe.Hostname)
				continue
			}

			find := true
			openReqsR2H, find = ListFindAndAdd(openReqsR2H,
							    fe.Hostname)
			if !find {
				go DoJoin(fe.Hostname, hierarchy, joinChan)
			}
		}
		//case <-time.After(1 * time.Second):
		//	//check len before waiting
	}
	}

	// same thing for openReqsH2R. but dont add in openReqsR2H
	for len(openReqsH2R) > 0 {
	select {
	case jr := <-joinChan:
		//send the join req, add it maybe?
		ctrlChan<-ctrlMsg{msgType:MSG_TYPE_RUN_GETFE,
				  hierarchy:hierarchy,
				  fe:*jr.Self}
		openReqsH2R = ListFindAndPop(openReqsH2R, jr.Self.Hostname)

		//check threshold start lower join cluster request if possible
		rtt_avg := (jr.RttMs + (Time64() - jr.Time)/NANO_TO_MILLISEC)
		if !lowerClusterFound &&
		   rtt_avg < m.config.Threshold[hierarchy] {
			go JoinCluster(m, hierarchy + 1, jr.Self.Hostname)
			lowerClusterFound = true
		}

		for _, fe := range jr.FEList {
			if !NodeIdBetween(*fe.Id, rootId, halfId) {
				continue
			}

			find := true
			openReqsH2R, find = ListFindAndAdd(openReqsH2R,
							    fe.Hostname)
			if !find {
				go DoJoin(fe.Hostname, hierarchy, joinChan)
			}
		}
	}
	}

	if !lowerClusterFound {
		// start own cluster. TODO
	}

	DLog("JoinCluster() complete")
	//signal completion of join
	ctrlChan<-ctrlMsg{msgType:MSG_TYPE_JOIN_COMPLETE, hierarchy:hierarchy}
}

func FillFingerTable(m *mcaster) {
	selfId := m.id
	selfFE := FingerEntry{Id:&selfId, Hostname:m.hostname}
	for i,_ := range m.topFT {
		m.topFT[i] = selfFE
	}
	for i,_ := range m.midFT {
		m.midFT[i] = selfFE
	}
	for i,_ := range m.botFT {
		m.botFT[i] = selfFE
	}
	if m.hostname == ucsdNode {
		HLog("The root node. return")
		// we are the first one. everything is set to selfFE
		return
	}

	// code to fill up the finger table
	joined := 0 // number of clusters joined
	for {
	select {
	case cmsg := <-ctrlChan:
		DLog("cmsg %v", cmsg)
		// receive messages from JoinCluster, add to Finger Tables
		switch cmsg.msgType {
		case MSG_TYPE_RUN_GETFE:  //0

		case MSG_TYPE_JOIN_COMPLETE : //2
			joined = joined + 1
			if joined == 3 {
				MLog("Joined all three clusters")
				break;
			}
		}

	}
	}
}

func DoFwdPkt(fe FingerEntry, limit NodeId, pkt FwdPacket, hierarchy int) {
	//call rpc to fwd it
	MLog("Enter DoFwdPkt fe %v, hierarchy %d", fe, hierarchy)

	conn, err := grpc.Dial(fe.Hostname, grpc.WithInsecure(), grpc.WithBlock());
	if err!= nil {
		HLog("Dial failed: fe %v", fe.Hostname)
		return
	}
	defer conn.Close()

	new_pkt := new(FwdPacket)
	hops := int32(len(pkt.EvalList) + 1)

	new_pkt.Limit = &limit;
	new_pkt.EvalList = append(new_pkt.EvalList, pkt.EvalList...)
	new_pkt.EvalList = append(pkt.EvalList,
				  &EvalInfo{Hops: hops,
					    Time: Time64(),
					    Node: &m.id,
				  })
	new_pkt.Payload = pkt.Payload
	new_pkt.Src = &m.id
	DLog("DoFwdPkt new_pkt: %v", new_pkt);

	c := NewMcasterClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(),
					   500 * time.Millisecond) //wait for 500ms
	defer cancel()
	_, err = c.Fwd(ctx, new_pkt)
	if err!= nil {
		HLog("Fwd to %v failed: %v", fe.Hostname, err)
		//send a message to ctrlChan, about the failure
		ctrlChan <- ctrlMsg{msgType:MSG_TYPE_FAILED_FE,
				    hierarchy:hierarchy,
				    fe:fe,
				   }
		return
	}
	DLog("FWD to %v complete", fe.Hostname)
}

func Time64() int64 {
	// returns nanoseconds since 1970
	// see https://golang.org/pkg/time/#Time.UnixNano
	// ... The result does not depend on the location ...
	ts := int64(time.Now().UnixNano())
	return ts
}

func NeedToFwd(fe FingerEntry, pkt FwdPacket, hierarchy int) bool {
	x:= __NeedToFwd(fe, pkt, hierarchy)
	DLog("NeedToFwd [%v], fe %v, pkt %v, hierarchy %d", x, fe, pkt, hierarchy)
	return x
}

func __NeedToFwd(fe FingerEntry, pkt FwdPacket, hierarchy int) bool {
	s := *pkt.Src
	//l := *pkt.Limit
	if fe.Id == nil { //unfilled fingerEntry
		return false
	}

	//ERROR: invalid indirect of pkt.Src.B (type int64) ?? WHATISTHIS
	if s.A == m.id.A && s.B == m.id.B && s.C == m.id.C {
		// we are are the source. fwd to everyone we know, except self
		if fe.Hostname == m.hostname {
			return false
		}
		return true
	}

	if s.A != m.id.A {
		// diff top level hierarchy:
		//      - fwd in top level within limit
		//      - fwd to all in bottom and lower levels
		if hierarchy == 1 && fe.Id.A < pkt.Limit.A {
			return true
		}
		if hierarchy > 1 {
			return true
		}

		return false
	} else if pkt.Src.B != m.id.B {
		// same top level, but diff mid level hierarchy:
		//      - fwd in mid level within limit
		//      - fwd to all in bottom and lower levels
		if hierarchy == 2 && fe.Id.B < pkt.Limit.B {
			return true
		}
		if hierarchy > 2 {
			return true
		}

		return false
	} else {
		// same top and mid level hierarchy:
		//  - fwd in bottom within limit
		if hierarchy == 3 && fe.Id.C < pkt.Limit.C {
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

	MLog("Configuration: %v", c)
	return c
}

func StartMcast() {
	for {
		DLog("Mcast start")
		//send out X bytes
		xbytes := make([]byte, m.config.XBytes)//zeros
		pkt := new(FwdPacket)
		pkt.Payload = append(pkt.Payload, xbytes...)
		self := m.id
		pkt.Src = &self
		pktChan <- *pkt

		// sleep for Y seconds
		time.Sleep(time.Duration(m.config.YSeconds) * time.Second)
		DLog("Mcast Complete")
	}
}

func ValidFE_ptr(felist []FingerEntry) []*FingerEntry {
	selfId := new(NodeId)
	*selfId = m.id
	selfFE := FingerEntry{Id:selfId, Hostname:m.hostname}

	retlist := make([]*FingerEntry, 0)
	for i := 0; i < len(felist); i++ {
		if !isEqualFE(felist[i], selfFE) { // add non invalid entries
			retlist = append(retlist, &felist[i])
		}
	}
	DLog("NonDupFE_ptr retlist: %v", retlist)
	return retlist
}

func ValidFE(felist []FingerEntry) []FingerEntry {
	selfId := new(NodeId)
	*selfId = m.id
	selfFE := FingerEntry{Id:selfId, Hostname:m.hostname}

	retlist := make([]FingerEntry, 0)
	for i := len(felist)-1; i >= 0; i-- { // return in reverse
		if !isEqualFE(felist[i], selfFE) { // add non invalid entries
			retlist = append(retlist, felist[i])
		}
	}
	DLog("ValidFE retlist: %v", retlist)
	return retlist
}

func addFE(hierarchy int, fe FingerEntry){
	switch hierarchy {
	case 1:
		__addFE(&m.topFT, fe) // TODO use Runze's func
	case 2:
		__addFE(&m.midFT, fe) // TODO use Runze's func
	case 3:
		__addFE(&m.botFT, fe) // TODO use Runze's func
	}
}

// fix this, dont use topFTsize
func __addFE( felist *[topFTsize]FingerEntry, fe FingerEntry) {
	//TODO use Runze's func instead
}

func main() {
	fmt.Println("--- Hierarchical P2P Mcast ---")
	log.SetFlags(log.Ldate | log.Ltime | log.LUTC | log.Lmicroseconds)

	m = new(mcaster)
	progArgs := os.Args
	if len(progArgs) != 3 {
		log.Fatalf("Insufficient Args. Usage: ./hp2pcast <config>.json <port>")
	} else {
		MLog("Args %v", progArgs[0])
	}
	progArgs = progArgs[1:]

	m.config = ReadConfig(progArgs[0])
	m.hostname = m.config.Hostname
	//m.config.UcsdNode = "localhost:50000"
	ucsdNode = m.config.UcsdNode
	m.s = grpc.NewServer()

	port := ":" + progArgs[1]
	m.hostname = m.hostname + port
	m.config.Hostname = m.hostname

	hash := McastHash(m.hostname)
	MLog("hostname %v: hash %v", m.hostname, hash)
	m.id = NodeId{A:hash, B:hash, C:hash}

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalln("failed to open tcp socket: %v", err)
	}
	RegisterMcasterServer(m.s, m)

	HLog("Start RPC server")
	go ServeRPCServer(m.s, lis)

	go JoinCluster(m, 1, ucsdNode)
	FillFingerTable(m)

	if m.config.MulticastFlag == 1 {
		go StartMcast()
	}

	for {
	select {

	case pkt := <-pktChan:
		MLog("received pkt Limit %v, Src %v", pkt.Limit, pkt.Src)
		DLog("received pkt %v", pkt)
		//decide fingertable based on m.id and pkt.src

		// get ValidFE in reverse. The limit is self for A/2 node,
		// for all others its the next node in the table.
		limit := m.id
		for _, fe := range ValidFE(m.topFT[:]) {
			if NeedToFwd(fe, pkt, 1) {
				go DoFwdPkt(fe, limit, pkt, 1)
			}
			limit = *fe.Id
		}

		limit = m.id
		for _, fe := range ValidFE(m.midFT[:]) {
			if NeedToFwd(fe, pkt, 2) {
				go DoFwdPkt(fe, limit, pkt, 2)
			}
			limit = *fe.Id
		}

		limit = m.id
		for _, fe := range ValidFE(m.botFT[:]) {
			if NeedToFwd(fe, pkt, 3) {
				go DoFwdPkt(fe, limit, pkt, 3)
			}
			limit = *fe.Id
		}

	// to handle failures
	case cmsg := <-ctrlChan:
		// TODO: handle failures
		switch cmsg.msgType {
		case MSG_TYPE_RUN_GETFE:
			//FIXME
		case MSG_TYPE_ADD_NODEID:
			DLog("MSG_TYPE_ADD_NODEID, add fe %v", cmsg.fe)
			addFE(cmsg.hierarchy, cmsg.fe)

		case MSG_TYPE_FAILED_FE:

			DLog("MSG_TYPE_ADD_NODEID, add fe %v", cmsg.fe)
			// 0. Invalidate the fe
			invalidateFE(cmsg.fe, cmsg.hierarchy)
			if cmsg.hierarchy == 3 {
				// nothing else to do
				continue
			}

			// 1. query in the lowest hierarchy abt fe
			fe, err := getRandomFE(3)
			// TODO: can we ask someone in higher levels?
			if err != nil {
				HLog("Error: %v", err)
				continue
			}

			limit := cmsg.fe.Id
			switch cmsg.hierarchy {
				case 1:
					limit.A = limit.A + 1 //TODO: modulo?
				case 2:
					limit.B = limit.B + 1
			}
			go DoGetFE(fe, cmsg.hierarchy, *limit)
			// 2. on recving the response replace old fe by sending
			//    a ctrlMsg{MSG_TYPE_ADD_NODEID}
		}
	}
	}
}
