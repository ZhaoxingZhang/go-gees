package main

import (
	"encoding/json"
	"fmt"
	"github.com/ZhaoxingZhang/geeorm/log"
	"github.com/ZhaoxingZhang/geerpc"
	"github.com/ZhaoxingZhang/geerpc/codec"
	"net"
	"sync"
	"time"
)

func startServer(addr chan string) {
	// pick a free port
	l, err := net.Listen("tcp", ":0")
	if err != nil {
		log.Errorf("network error:", err)
	}
	log.Infof("start rpc server on %v", l.Addr())
	addr <- l.Addr().String()
	geerpc.Accept(l)
}

func simpleClient(addr chan string) {
	// in fact, following code is like a simple geerpc client
	conn, _ := net.Dial("tcp", <-addr)
	defer func() { _ = conn.Close() }()

	time.Sleep(time.Second)
	// send options
	_ = json.NewEncoder(conn).Encode(geerpc.DefaultOption)
	cc := codec.NewGobCodec(conn)
	// send request & receive response
	for i := 0; i < 5; i++ {
		h := &codec.Header{
			ServiceMethod: "Foo.Sum",
			Seq:           uint64(i),
		}
		_ = cc.Write(h, fmt.Sprintf("geerpc req %d", h.Seq))
		_ = cc.ReadHeader(h)
		var reply string
		_ = cc.ReadBody(&reply)
		log.Infof("reply: %v", reply)
	}
}
func clientSDK(addr chan string) {
	var wg sync.WaitGroup
	client, err := geerpc.Dial("tcp", <-addr)
	defer client.Close()
	if err != nil {
		log.Errorf("Dial error: %v", err)
	}
	//time.Sleep(time.Second)
	serviceMethod := "Foo.Sum"
	done := make(chan *geerpc.Call, 3)
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(i int) {
			args := fmt.Sprintf("geerpc req %d", i)
			defer wg.Done()
			var reply string
			call := client.Go(serviceMethod, args, &reply, done)
			if _, ok := <-call.Done; ok {
				log.Infof("reply: %v", reply)
			}

		}(i)
	}
	wg.Wait()

}
func main() {
	addr := make(chan string)
	go startServer(addr)

	//simpleClient(addr)
	clientSDK(addr)
}
