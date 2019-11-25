package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/pprof"

	proto "github.com/micro/examples/stream/server/proto"
	"github.com/micro/go-micro/service/grpc"
)

func bidirectional(cl proto.StreamerService) {
	// create streaming client
	stream, err := cl.Stream(context.Background())
	if err != nil {
		fmt.Println("err:", err)
		return
	}

	// bidirectional stream
	// send and receive messages for a 10 count
	for j := 0; j < 10; j++ {
		if err := stream.Send(&proto.Request{Count: int64(j)}); err != nil {
			fmt.Println("err:", err)
			return
		}
		rsp, err := stream.Recv()
		if err != nil {
			fmt.Println("recv err", err)
			break
		}
		fmt.Printf("Sent msg %v got msg %v\n", j, rsp.Count)
	}

	// close the stream
	if err := stream.Close(); err != nil {
		fmt.Println("stream close err:", err)
	}
}

func serverStream(cl proto.StreamerService) {
	// send request to stream count of 10
	stream, err := cl.ServerStream(context.Background(), &proto.Request{Count: int64(10)})
	if err != nil {
		fmt.Println("err:", err)
		return
	}

	// server side stream
	// receive messages for a 10 count
	for j := 0; j < 10; j++ {
		rsp, err := stream.Recv()
		if err != nil {
			fmt.Println("recv err", err)
			break
		}
		fmt.Printf("got msg %v\n", rsp.Count)
	}

	// close the stream
	if err := stream.Close(); err != nil {
		fmt.Println("stream close err:", err)
	}
}

func main() {
	go func() { panic(profilerSetup()) }()

	service := grpc.NewService()
	service.Init()

	// create client
	cl := proto.NewStreamerService("go.micro.srv.stream", service.Client())

	// bidirectional stream
	// bidirectional(cl)

	// server side stream
	for {
		serverStream(cl)
	}
}

func profilerSetup() error {
	mux := &http.ServeMux{}

	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
	mux.HandleFunc("/debug/pprof/block", pprof.Handler("block").ServeHTTP)
	mux.HandleFunc("/debug/pprof/heap", pprof.Handler("heap").ServeHTTP)
	mux.HandleFunc("/debug/pprof/goroutine", pprof.Handler("goroutine").ServeHTTP)
	mux.HandleFunc("/debug/pprof/threadcreate", pprof.Handler("threadcreate").ServeHTTP)

	if err := http.ListenAndServe(":6060", mux); err != nil {
		return err
	}

	return errors.New("pprof server stopped")
}
