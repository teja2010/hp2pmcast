package main

import (
	"os"
	//"io"
	"net"
	"fmt"
	"log"
	"time"
	"errors"
	"context"
	"runtime"
	"strconv"
	"encoding/json"

	"google.golang.org/grpc"
)

var (
	n *narada
	pktChan chan NFwdPacket
	logChan chan NFwdPacket
	initState = true
)

func Time64() int64 {
	// returns nanoseconds since 1970
	// see https://golang.org/pkg/time/#Time.UnixNano
	// ... The result does not depend on the location ...
	ts := int64(time.Now().UnixNano())
	return ts
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
	if n.config.Dlog == 0 {
		return
	}

	file := LogFormat()
	if n.config.LogColor == 1{
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
	if n.config.LogColor == 1{
		log.Printf("\x1B[31m HIGH  | " + file+ format + "\033[0m", v ...)
	} else {
		log.Printf(" HIGH  | " +file+ format , v ...)
	}
}

type ntwConfig struct {
	Narada [][]string
}

func fillNtwfromConfig(NtwConfig string) {
	fp, err := os.Open(NtwConfig)
	if err != nil {
		log.Fatalf("Config file open failed %v", err)
	}
	defer fp.Close()

	decoder := json.NewDecoder(fp)
	ntw := ntwConfig{}
	err = decoder.Decode(&ntw)
	if err != nil {
		log.Fatalf("Config file read failed: %v", err)
	}

	MLog("NtwConfig: %v, file %v", ntw, NtwConfig)

	for _, ll := range ntw.Narada {
		if n.hostname == ll[0] {
			n.links = append(n.links, ll[1:]...)
		}
	}

	MLog("Links %v", n.links)
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

type Configuration struct {
	Hostname string
	NtwConfig string
	MulticastFlag int

	XBytes int
	YSeconds int
	Dlog int
	LogColor int
}

type narada struct {
	UnimplementedNaradamcastServer
	s *grpc.Server

	hostname string
	config Configuration
	links []string
}

func StartMcast() {
	seqNum := 0
	for {
		// sleep for Y seconds
		time.Sleep(time.Duration(n.config.YSeconds) * time.Second)

		DLog("Mcast start")
		//send out X bytes
		xbytes := make([]byte, n.config.XBytes)//zeros
		pkt := new(NFwdPacket)
		pkt.Payload = append(pkt.Payload, xbytes...)
		pkt.SrcHostname = n.hostname
		hops := int32(len(pkt.EvalList) + 1)
		pkt.EvalList = append(pkt.EvalList,
					  &NEvalInfo{Hops: hops,
						    Time: Time64(),
						    Hostname: n.hostname,
					  })
		pkt.SeqNum = int32(seqNum)
		seqNum = seqNum + 1
		pktChan<- *pkt
		logChan<-*pkt
		DLog("Mcast Complete")

	}
}

func sendToLogger() {
	type logconfig struct {
		NLoggerHostname string
	}

	fp, err := os.Open("logger.json")
	if err != nil {
		log.Fatalf("Config file open failed %v", err)
	}
	decoder := json.NewDecoder(fp)
	lc := logconfig{}
	err = decoder.Decode(&lc)
	if err != nil {
		log.Fatalf("Config file read failed: %v", err)
	}
	defer fp.Close()

	conn, err := grpc.Dial(lc.NLoggerHostname, grpc.WithInsecure());
	if err!= nil {
		HLog("Dial failed: logger %v", lc.NLoggerHostname)
	}

	for {
	select {
	case pkt := <-logChan:

		MLog("received pkt %v", pkt)

		new_pkt := new(NFwdPacket)
		new_pkt.EvalList = append(new_pkt.EvalList, pkt.EvalList...)
		new_pkt.SeqNum = pkt.SeqNum
		new_pkt.SrcHostname = n.hostname
		c := NewNaradamcastClient(conn)
		ctx, cancel := context.WithTimeout(context.Background(),
						   1 * time.Second) //wait for 1s
		defer cancel()
		_, err = c.NaradaFwd(ctx, new_pkt)
		if err!= nil {
			HLog("Fwd to %v failed: %v", lc.NLoggerHostname, err)
			//send a message to ctrlChan, about the failure
			conn.Close()
			conn, err = grpc.Dial(lc.NLoggerHostname, grpc.WithInsecure());
			if err!= nil {
				HLog("Dial failed: logger %v", lc.NLoggerHostname)
			}
		}
	}}

}

func (s *narada) NaradaFwd(ctx context.Context, in *NFwdPacket) (*NEmpty, error) {
	DLog("Fwd in: NFwdPacket %v", in)
	if initState {
		return nil, errors.New("Initializing...")
	}
	out := new(NEmpty)

	// todo log the arrival of the packet, send it out periodically
	hops := int32(len(in.EvalList) + 1)
	in.EvalList = append(in.EvalList,
				  &NEvalInfo{Hops: hops,
					    Time: Time64(),
					    Hostname: n.hostname,
				  })

	// send to pktChan, it can decide to forward the pkt to others
	pktChan<-*in
	logChan<-*in

	return out, nil
}

func DoFwdPkt(hostname string, pkt NFwdPacket) {
	//call rpc to fwd it
	MLog("Enter DoFwdPkt hostname %v", hostname)

	conn, err := grpc.Dial(hostname, grpc.WithInsecure());
	if err!= nil {
		HLog("Dial failed: fe %v", hostname)
		return
	}
	defer conn.Close()

	new_pkt := new(NFwdPacket)
	hops := int32(len(pkt.EvalList) + 1)

	new_pkt.SrcHostname = n.hostname;
	new_pkt.SeqNum = pkt.SeqNum
	new_pkt.EvalList = append(new_pkt.EvalList, pkt.EvalList...)
	new_pkt.EvalList = append(pkt.EvalList,
				  &NEvalInfo{Hops: hops,
					    Time: Time64(),
					    Hostname: n.hostname,
				  })
	new_pkt.Payload = pkt.Payload
	DLog("DoFwdPkt new_pkt: %v", new_pkt);

	c := NewNaradamcastClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(),
					   500 * time.Millisecond) //wait for 500ms
	defer cancel()
	_, err = c.NaradaFwd(ctx, new_pkt)
	if err!= nil {
		HLog("Fwd to %v failed: %v", hostname, err)
		return
	}
	DLog("FWD to %v complete", hostname)
}

func ServeRPCServer (s *grpc.Server, lis net.Listener) {
	err := s.Serve(lis);
	if err != nil {
		log.Fatalf("Failed to start server %v", err)
	}
	DLog("Started RPC Server")
}

func main() {
	fmt.Println("--- Narada Fwding ---")
	log.SetFlags(log.Ldate | log.Ltime | log.LUTC | log.Lmicroseconds)

	n = new(narada)
	progArgs := os.Args
	if len(progArgs) != 3 {
		log.Fatalf("Insufficient Args. Usage: ./narada <config>.json <port>")
	} else {
		MLog("Args %v", progArgs[0])
	}
	progArgs = progArgs[1:]

	n.config = ReadConfig(progArgs[0])

	n.hostname = n.config.Hostname
	port := ":" + progArgs[1]
	n.hostname = n.hostname + port
	n.config.Hostname = n.hostname
	DLog("Hostname %v", n.hostname)

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalln("failed to open tcp socket: %v", err)
	}
	n.s = grpc.NewServer()
	RegisterNaradamcastServer(n.s, n)
	pktChan = make(chan NFwdPacket, 100)
	logChan = make(chan NFwdPacket, 100)
	go sendToLogger()

	HLog("Start RPC server")
	go ServeRPCServer(n.s, lis)

	fillNtwfromConfig(n.config.NtwConfig)

	initState = false
	if n.config.MulticastFlag == 1 {
		go StartMcast()
	}

	DLog("Waiting for packets")
	for {
	select {

	case pkt := <-pktChan:
		MLog("received pkt Src %v", pkt.SrcHostname)
		DLog("received pkt %v", pkt)

		for _, ll := range n.links {
			if ll != pkt.SrcHostname {
				go DoFwdPkt(ll, pkt)
			}
		}


	}
	}
}
