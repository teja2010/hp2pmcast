##Hierarchical Peer2Peer Multicast

### requirements:
`sudo apt  install protobuf-compiler golang-go`

`go get -u google.golang.org/grpc`

`go get -u github.com/golang/protobuf/protoc-gen-go`

add this to ~/.profile:

    # set PATH so it includes user's go bin
    if [ -d "$HOME/go/bin" ] ; then
        PATH="$HOME/go/bin:$PATH"
    fi



### Build instructions:

    ./codegen mcast.proto
    make

### Run
    ./hp2pmcast <config>.json <port>

