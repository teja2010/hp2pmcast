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
	FTsize = 64
	CHORD_HIERARCHY = 0 // all non-zero rings are higher clusters
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
	Hierarchies int
	XBytes int
	YSeconds int
	//Add new items here. read from config.json
	Threshold []int64
	UcsdNode string
	Dlog int
	LogColor int
}
func (c Configuration) String() string {
	return fmt.Sprintf("{Hostname %v, MulticastFlag %v, Hierarchies %d " +
			   "XBytes %v, YSeconds %v, Threshold %v, UcsdNode %v}",
			   c.Hostname, c.MulticastFlag, c.Hierarchies,
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
	msgTypeStr := ""
	switch c.msgType {
	case MSG_TYPE_RUN_GETFE:
		msgTypeStr = "MSG_TYPE_RUN_GETFE"
	case MSG_TYPE_ADD_NODEID:
		msgTypeStr = "MSG_TYPE_ADD_NODEID"
	case MSG_TYPE_JOIN_COMPLETE:
		msgTypeStr = "MSG_TYPE_JOIN_COMPLETE"
	case MSG_TYPE_FAILED_FE:
		msgTypeStr = "MSG_TYPE_FAILED_FE"
	}
	return fmt.Sprintf("{ msgType %v, hierarchy %v, " +
			   "FingerEntry %v}",
			   msgTypeStr, c.hierarchy, c.fe)
}

type mcaster struct {
	UnimplementedMcasterServer

	s *grpc.Server
	id NodeId
	hostname string
	FT [][]FingerEntry
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
	out.FEList = ValidFE_ptr(m.FT[in.Hierarchy])
	DLog("JoinResp %v", out)

	return out, nil
}

func (s *mcaster) GetFingerEntry(ctx context.Context, in *GetFERequest) (*GetFEResponse, error) {
	DLog("GetFingerEntry : GetFERequest %v", in)
	out := new(GetFEResponse)
	out.Hierarchy = in.Hierarchy
	out.Src = &m.id

	return out, nil
}

func (s *mcaster) SetSuccessor(ctx context.Context, in *Successor) (*Empty, error) {
	DLog("SetSuccessor Id %v", in.FE.Id)
	out := new(Empty)

	return out, nil
}

// call this is for any receiving any message
func isValidNodeId(n NodeId) bool {
	return (len(n.Ids) == m.config.Hierarchies)
}


func isEqualNodeId(n1, n2 NodeId) bool {
	if !isValidNodeId(n1) || !isValidNodeId(n2) {
		return false
	}

	for i:=0; i< m.config.Hierarchies; i++ {
		if n1.Ids[i] != n2.Ids[i] {
			return false
		}
	}
	return true
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

	for i := 0; i < FTsize; i++ {
		if isEqualFE(m.FT[hierarchy][i], fe) {
			m.FT[hierarchy][i] = selfFE
		}
	}
}

func getRandomFE(hierarchy int) (FingerEntry, error) {
	selfId := new(NodeId)
	*selfId = m.id
	selfFE := FingerEntry{Id:selfId, Hostname:m.hostname}

	rr := rand.Intn(FTsize)
	for i := 0; i < FTsize; i++ {
		if isEqualFE(m.FT[hierarchy][(rr + i)%FTsize], selfFE) {
			return m.FT[hierarchy][(rr + i)%FTsize], nil
		}
	}

	return selfFE, errors.New("Unable to find a rand FE")
}

func DoGetFE(dest string, hierarchy int, limit NodeId) {
	// send a req to fe, @ hierarchy, get an entry within limit

	MLog("Enter DoGetFE dest %v, hierarchy %d, limit %v",
		dest, hierarchy, limit)

	conn, err := grpc.Dial(dest, grpc.WithInsecure(), grpc.WithBlock());
	if err!= nil {
		HLog("Dial failed: dest %v", dest)
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
		HLog("GetFingerEntry to %v failed, err %v", dest, err)
		// TODO what to do?
		return
	}
	DLog("GetFingerEntry complete:  %v", resp)

	if resp.Src == nil || resp.NewFE == nil {
		HLog("Incomplete response")
		return
	}

	if isEqualNodeId(*resp.Src, *resp.NewFE.Id) {
		// this is the entry that should be added.
		ctrlChan <- ctrlMsg{msgType:MSG_TYPE_ADD_NODEID,
				    hierarchy : hierarchy,
				    fe : *resp.NewFE,
				   }
	} else {
		// send another GetFE request.
		ctrlChan <- ctrlMsg{msgType:MSG_TYPE_RUN_GETFE,
				    hierarchy : hierarchy,
				    fe : *resp.NewFE,
				   }
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

func HalfId(n NodeId) NodeId {
	var half NodeId
	topBit := uint64(1) << (FTsize-1);
	for h:=0; h<m.config.Hierarchies; h++ {
		half.Ids = append(half.Ids, n.Ids[h] ^ topBit)
	}

	return half
}

func NodeIdBetween(id, left, right NodeId, hierarchy int) bool {
	ret := __NodeIdBetween(id, left, right, hierarchy)
	DLog("NodeIdBetween[%v]: left %v < id %v < right %v",
		ret, left, id, right)
	return ret
}

func __NodeIdBetween(id, left, right NodeId, hierarchy int) bool {
	// is left < id < right ?
	if left.Ids[hierarchy] < right.Ids[hierarchy] {
		// -----|------|-------|-----
		//     left < id  <  right
		return ( left.Ids[hierarchy] < id.Ids[hierarchy] &&
		        right.Ids[hierarchy] > id.Ids[hierarchy] )
	} else if left.Ids[hierarchy] > right.Ids[hierarchy] {
		// -----|------|-------|-----  OR  -----|-------|-------|-----
		//     id  <  right < left            right  > left >  id
		return ((  left.Ids[hierarchy] > id.Ids[hierarchy] &&
			  right.Ids[hierarchy] > id.Ids[hierarchy] ) ||
		        (  left.Ids[hierarchy] < id.Ids[hierarchy] &&
			  right.Ids[hierarchy] < id.Ids[hierarchy] ))
	}

	return false
}

func JoinChordRing(rootNode string) {
	HLog("JoinChordRing rootNode %v", rootNode)

	// for each empty entry send a getFE request
	//    send getFE requests based on responses.
	//    if dest_host == response, fwd response to CtrlChan so it adds entry

	for i:=uint64(0); i< FTsize; i++ {
		var limitId NodeId
		limitId.Ids = append(limitId.Ids, m.id.Ids...)
		limitId.Ids[CHORD_HIERARCHY] = (limitId.Ids[CHORD_HIERARCHY] +
						 (uint64(1) <<i) )
		go DoGetFE(rootNode, CHORD_HIERARCHY, limitId)
	}
}

func JoinCluster(hierarchy int, rootNode string) {
	if hierarchy < CHORD_HIERARCHY {
		log.Fatalf("hierarchy less than CHORD_HIERARCHY")
	} else if hierarchy == CHORD_HIERARCHY {
		// Run simple chord here.
		JoinChordRing(rootNode)
		return
	}
	HLog("JoinCluster hierarchy %d, rootNode %s", hierarchy, rootNode)

	hash := McastHash(rootNode)
	var rootId NodeId
	for h:=0; h<m.config.Hierarchies; h++ {
		rootId.Ids = append(rootId.Ids, hash)
	}

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
			go JoinCluster(hierarchy -1, jr.Self.Hostname)
			lowerClusterFound = true
		}

		for _, fe := range jr.FEList {
			if !NodeIdBetween(*fe.Id, rootId, halfId, hierarchy) {
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
			go JoinCluster(hierarchy - 1, jr.Self.Hostname)
			lowerClusterFound = true
		}

		for _, fe := range jr.FEList {
			if !NodeIdBetween(*fe.Id, rootId, halfId, hierarchy) {
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
	var selfId NodeId
	selfId.Ids = append(selfId.Ids, m.id.Ids...)
	selfFE := FingerEntry{Id:&selfId, Hostname:m.hostname}
	for h:=0 ; h < m.config.Hierarchies; h++ {
		for i,_ := range m.FT[h] {
			m.FT[h][i] = selfFE // set everything to invalid FE
		}
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

		case MSG_TYPE_ADD_NODEID: //1

		case MSG_TYPE_JOIN_COMPLETE : //2
			joined = joined + 1
			if joined == m.config.Hierarchies {
				MLog("Joined all %d clusters", joined)
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
	if isEqualNodeId(s, m.id) {
		// we are are the source. fwd to everyone we know, except self
		if fe.Hostname == m.hostname {
			return false
		}
		return true
	}

	for h:=0 ; h<m.config.Hierarchies; h++ {
		if s.Ids[h] != m.id.Ids[h] {
			// diff h level hierarchy
			//    - fwd at h level within limit
			if hierarchy == h && fe.Id.Ids[h] < pkt.Limit.Ids[h] {
				return true
			}

			//    - fwd to all lower levels
			if hierarchy < h {
				return true
			}

			//    - dont forward to higher levels and over limit
			return false
		}
	}

	// in the previous loop, any difference in hierarchy id will cause a
	// return. If it arrives here, that means all the ids match selfId.
	return false
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

func addFE(hierarchy int, fe FingerEntry) {
	//__addFE(&m.FT[hierarchy], fe) // TODO use Runze's func
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

	//create and set all FT to self

	m.hostname = m.config.Hostname
	//m.config.UcsdNode = "localhost:50000"
	ucsdNode = m.config.UcsdNode
	m.s = grpc.NewServer()

	port := ":" + progArgs[1]
	m.hostname = m.hostname + port
	m.config.Hostname = m.hostname

	hash := McastHash(m.hostname)
	MLog("hostname %v: hash %v", m.hostname, hash)
	for h:=0; h<m.config.Hierarchies; h++ {
		m.id.Ids = append(m.id.Ids, hash)
	}

	//fill up FT with selfId
	selfId := new(NodeId)
	*selfId = m.id
	selfFE := FingerEntry{Id:selfId, Hostname:m.hostname}
	for h:=0; h<m.config.Hierarchies; h++ {
		var hEntry []FingerEntry
		m.FT = append(m.FT, hEntry)
		for i:=0; i<FTsize; i++ {
			m.FT[h] = append(m.FT[h], selfFE)
		}
	}

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalln("failed to open tcp socket: %v", err)
	}
	RegisterMcasterServer(m.s, m)

	HLog("Start RPC server")
	go ServeRPCServer(m.s, lis)

	
	go JoinCluster(m.config.Hierarchies -1, ucsdNode)
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
		for h:=0; h<m.config.Hierarchies; h++ {
			limit := m.id
			for _, fe := range ValidFE(m.FT[h]) {
				if NeedToFwd(fe, pkt, h) {
					go DoFwdPkt(fe, limit, pkt, h)
				}
				limit = *fe.Id
			}
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
			if cmsg.hierarchy == CHORD_HIERARCHY {
				// nothing else to do
				continue
			}

			// 1. query in the lowest hierarchy abt fe
			fe, err := getRandomFE(CHORD_HIERARCHY)
			// TODO: can we ask someone in higher levels?
			if err != nil {
				HLog("Error: %v", err)
				continue
			}

			var limit NodeId
			limit = *cmsg.fe.Id
			limit.Ids[cmsg.hierarchy]++ // TODO: modulo?

			go DoGetFE(fe.Hostname, cmsg.hierarchy, limit)
			// 2. on recving the response replace old fe by sending
			//    a ctrlMsg{MSG_TYPE_ADD_NODEID}
		}
	}
	}
}
