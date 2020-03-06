package main

import (
	"os"
	"fmt"
	"log"
	"net"
	"errors"
	"runtime"
	"strconv"
	"strings"
	"context"
	//"io/ioutil"
	"encoding/json"

	"google.golang.org/grpc"
)

type logStr struct {
	hostname string
	log string
}

var (
	asd = ""
	logChan = make(chan logStr, 100)
)

func fileLog(seq int32, hostname string, evList []*EvalInfo) {
	//find the file, and write to it
	var ss string
	ss = fmt.Sprint(seq)
	for _, e := range(evList) {
		ss += ";" + fmt.Sprint(e.Time) + "," + e.Hostname
	}
	ss += "\n"

	logChan<-logStr{hostname:hostname, log:ss}
}
func fileLog2(seq int32, hostname string, evList []*NEvalInfo) {
	//find the file, and write to it
	var ss string
	ss = fmt.Sprint(seq)
	for _, e := range(evList) {
		ss += ";" + fmt.Sprint(e.Time) + "," + e.Hostname
	}
	ss += "\n"
	MLog("%v",ss)

	logChan<-logStr{hostname:hostname, log:ss}
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

func MLog(format string, v ...interface{}) {
	file := LogFormat()
	log.Printf(" MID   | " + file + format , v ...)
}

func (s *logger) NaradaFwd(ctx context.Context, in *NFwdPacket) (*NEmpty, error) {
	out := new(NEmpty)

	go fileLog2(in.SeqNum, in.SrcHostname, in.EvalList)

	return out, nil
}

func (s *logger) Fwd(ctx context.Context, in *FwdPacket) (*Empty, error) {
	out := new(Empty)

	//MLog("in %v", in)
	go fileLog(in.SeqNum, in.SrcHostname, in.EvalList)

	return out, nil
}

func (s *logger) Join(ctx context.Context, in *JoinReq) (*JoinResp, error) {
	MLog("Join : JoinReq %v", in)
	return nil, errors.New("Not supported")
}
func (s *logger) GetFingerEntry(ctx context.Context, in *GetFERequest) (*GetFEResponse, error) {
	return nil, errors.New("Not supported")
}
func (s *logger) SetSuccessor(ctx context.Context, in *Successor) (*Empty, error) {
	return nil, errors.New("Not supported")
}
func (s *logger) SetPredecessor(ctx context.Context, in *Predecessor) (*Empty, error) {
	return nil, errors.New("Not supported")
}

type logger struct {
	UnimplementedMcasterServer
	UnimplementedNaradamcastServer

	s1, s2 *grpc.Server
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
	LoggerHostname string
	NLoggerHostname string
}
func startRPC1(l *logger, lis1 net.Listener) {
	l.s1 = grpc.NewServer()
	RegisterMcasterServer(l.s1, l)

	MLog("Started RPC Server1")
	err := l.s1.Serve(lis1);
	if err != nil {
		log.Fatalf("Failed to start server %v", err)
	}
}

func startRPC2(l *logger, lis2 net.Listener) {
	l.s2 = grpc.NewServer()
	RegisterNaradamcastServer(l.s2, l)

	MLog("Started RPC Server2")
	err := l.s2.Serve(lis2);
	if err != nil {
		log.Fatalf("Failed to start server %v", err)
	}
}

func openFile(filename string) *os.File {
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}

	return f
}

func main() {
	fmt.Println("--- Hierarchical P2P Mcast ---")
	log.SetFlags(log.Ldate | log.Ltime | log.LUTC | log.Lmicroseconds)
	l := new(logger)


	config := ReadConfig("logger.json")
	port := strings.Split(config.LoggerHostname, ":")

	lis1, err := net.Listen("tcp", ":" + port[1])
	if err != nil {
		log.Fatalln("failed to open tcp socket: %v", err)
	}
	go startRPC1(l, lis1)

	Nport := strings.Split(config.NLoggerHostname, ":")
	lis2, err := net.Listen("tcp", ":" + Nport[1])
	if err != nil {
		log.Fatalln("failed to open tcp socket: %v", err)
	}
	go startRPC2(l, lis2)


	fileMap := make(map[string]*os.File)

	for {
	select {
	case ls:= <-logChan :
		MLog("%v", ls)
		key := ls.hostname + ".log"
		fp, prs := fileMap[key]
		if !prs  {
			fp = openFile(key)
			fileMap[key] = fp
		}

		_, err := fp.WriteString(ls.log)
		if err != nil {
			MLog("WriteString failed %v", err)
		}
	}
	}
}

