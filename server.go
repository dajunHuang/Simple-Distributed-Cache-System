package main

import (
    "os"
	"fmt"
    "net"
    "log"
	"context"
    "regexp"
	"net/http"
	"io/ioutil"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
    pb  "djhuang.top/cacheserver/cache"
)

var server cacheServer
var address [4]string
var client[2] pb.CacheClient
var conn[2] *grpc.ClientConn

func setAddress() {
    if os.Args[1] == "1" {
        address[0] = "127.0.0.1:9527"
        address[1] = "127.0.0.1:9530"

        address[2] = "127.0.0.1:9531"
        address[3] = "127.0.0.1:9532"
    } else if os.Args[1] == "2" {
        address[0] = "127.0.0.1:9528"
        address[1] = "127.0.0.1:9531"

        address[2] = "127.0.0.1:9530"
        address[3] = "127.0.0.1:9532"
    } else if os.Args[1] == "3" {
        address[0] = "127.0.0.1:9529"
        address[1] = "127.0.0.1:9532"

        address[2] = "127.0.0.1:9530"
        address[3] = "127.0.0.1:9531"
    } else {
        fmt.Println("only 3 cacheserver.")
    }
}

func handleGet(w http.ResponseWriter, key string) {
    fmt.Println("get", key)

    if _, ok := server.cache[key]; ok {
	    w.WriteHeader(http.StatusOK)
        w.Header().Set("Content-Type", "application/json")
        fmt.Fprintln(w, "{\""+key+"\":\""+server.cache[key]+"\"}")
        return
    }

    w.WriteHeader(http.StatusNotFound)
}

func handleSet(w http.ResponseWriter, jsonstr string) {

    reg := regexp.MustCompile(`{\s*"(.*)"\s*:\s*"(.*)"\s*}`)
    if reg == nil { 
        fmt.Println("regexp err")
        return
    }
    result := reg.FindAllStringSubmatch(jsonstr, -1)
    key, value:= result[0][1], result[0][2]

    fmt.Println("set", key, ":", value)

    server.cache[key] = value
    CacheSet(client[0], &pb.SetRequest{Key:key, Value:value})
    CacheSet(client[1], &pb.SetRequest{Key:key, Value:value})

	w.WriteHeader(http.StatusOK)
}

func handleDelete(w http.ResponseWriter, key string) {
    fmt.Println("delete", key)
    if _, ok:= server.cache[key]; ok {
        delete(server.cache, key)
        CacheDelete(client[0], &pb.DeleteRequest{Key:key})
        CacheDelete(client[1], &pb.DeleteRequest{Key:key})

	    w.WriteHeader(http.StatusOK)
        fmt.Fprintln(w, "1")
        return
    }

	w.WriteHeader(http.StatusOK)
    fmt.Fprintln(w, "0")
}

func handleHttpRequest(w http.ResponseWriter, r *http.Request) {
    if r.Method == http.MethodGet {
        handleGet(w, r.URL.String()[1:])
    } else if r.Method == http.MethodPost {
        body, err := ioutil.ReadAll(r.Body)
        if err != nil {
            http.Error(w, "Unable to read request body.", http.StatusInternalServerError)
            return
        }
        handleSet(w, string(body))
    } else if r.Method == http.MethodDelete {
        handleDelete(w, r.URL.String()[1:])
    } else {
        http.Error(w, "Unsupport http request.", http.StatusMethodNotAllowed)
    }
}

type cacheServer struct {
    pb.UnimplementedCacheServer
    cache map[string]string
}

func (s *cacheServer) GetCache (ctx context.Context, req *pb.GetRequest) (*pb.GetReply, error) {
    return &pb.GetReply{Key:req.Key, Value:s.cache[req.Key]}, nil
}

func (s *cacheServer) SetCache (ctx context.Context, req *pb.SetRequest) (*pb.SetReply, error) {
    s.cache[req.Key] = req.Value
    return &pb.SetReply{}, nil
}

func (s *cacheServer) DeleteCache (ctx context.Context, req *pb.DeleteRequest) (*pb.DeleteReply, error) {
    if _, ok:= s.cache[req.Key]; ok {
        delete(s.cache, req.Key)
        return &pb.DeleteReply{Num: 1}, nil
    }
    return &pb.DeleteReply{Num: 0}, nil
}

func startHttpServer() {
    http.HandleFunc("/", handleHttpRequest)
    fmt.Println("Listening http on", address[0])
    err := http.ListenAndServe(address[0], nil)
    if err != nil {
        fmt.Println("Listten failed:", err)
    }
}

func startRpcServer() {
    fmt.Println("Listening rpc on", address[1])
    lis, err := net.Listen("tcp", address[1])
    if err != nil {
        log.Fatalf("failed to listen: %v", err)
    }

    var opts []grpc.ServerOption
    grpcServer := grpc.NewServer(opts...)
    server = cacheServer{cache: make(map[string]string)}
    pb.RegisterCacheServer(grpcServer, &server)
    grpcServer.Serve(lis)
}

func setupClient() {
	var opts []grpc.DialOption
    var err error
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	conn[0], err = grpc.Dial(address[2], opts...)
	if err != nil {
		fmt.Println("fail to dial: %v", err)
	}
    fmt.Println("Set up client for",address[2])

	conn[1], err = grpc.Dial(address[3], opts...)
	if err != nil {
		fmt.Println("fail to dial: %v", err)
	}
    fmt.Println("Set up client for",address[3])

    client[0] = pb.NewCacheClient(conn[0])
    client[1] = pb.NewCacheClient(conn[1])
}

func main() {
    if len(os.Args) != 2 {
        fmt.Println("please specify server index(1-3).")
        return
    }

    setAddress()
    go startHttpServer()
    go startRpcServer()
    setupClient()

    select{}
}
