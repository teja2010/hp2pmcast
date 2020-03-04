
all:
	go build -o hp2pmcast main.go mcast.pb.go
	go build -o narada narada.go narada.pb.go
