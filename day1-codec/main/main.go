package main

import (
	"fmt"
	"geerpc"
	"log"
	"net"
	"sync"
	"time"
)

func startServer(addr chan string) {
	// pick a free port
	l, err := net.Listen("tcp", ":0")
	if err != nil {
		log.Fatal("network error:", err)
	}
	log.Println("start rpc server on", l.Addr())
	addr <- l.Addr().String()
	geerpc.Accept(l)
}

func main() {
	addr := make(chan string)
	go startServer(addr)

	// in fact, following code is like a simple geerpc client
	client, _ := geerpc.Dial("tcp", <-addr)
	defer func() { _ = client.Close() }()

	time.Sleep(time.Second)
	// // send options
	// _ = json.NewEncoder(conn).Encode(geerpc.DefaultOption)
	// cc := codec.NewGobCodec(conn)
	// // send request & receive response
	// for i := 0; i < 5; i++ {
	// 	h := &codec.Header{
	// 		ServiceMethod: "Foo.Sum",
	// 		Seq:           uint64(i),
	// 	}
	// 	_ = cc.Write(h, fmt.Sprintf("geerpc req %d", h.Seq)) //send the request
	// 	_ = cc.ReadHeader(h)                                 //read the response header from server
	// 	var reply string
	// 	_ = cc.ReadBody(&reply) // read the response body from server
	// 	log.Println("reply:", reply)
	// }

	// send requst & receive response
	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			args := fmt.Sprintf("geerpc req %d", i)
			var reply string
			if err := client.Call("Foo.Sum", args, &reply); err != nil {
				log.Fatal("call Foo.Sum error:", err)
			}
			log.Println("reply:", reply)
		}(i)
	}
	wg.Wait()
}
