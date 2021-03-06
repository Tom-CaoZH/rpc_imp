package main

import (
	"geerpc"
	"log"
	"net"
	"sync"
	"time"
)

// every type is a service , the methods every type contains are the service.method
type Foo int

type Args struct{ Num1, Num2 int }

func (f Foo) Sum(args Args, reply *int) error {
	*reply = args.Num1 + args.Num2
	return nil
}

func (f Foo) Minus(args Args, reply *int) error {
	*reply = args.Num2 - args.Num1
	return nil
}

func startServer(addr chan string) {
	var foo Foo
	if err := geerpc.Register(&foo); err != nil {
		log.Fatal("register error:", err)
	}
	l, err := net.Listen("tcp", ":0")
	if err != nil {
		log.Fatal("network error:", err)
	}
	log.Println("start rpc server on", l.Addr())
	addr <- l.Addr().String()
	geerpc.Accept(l)
}

func main() {
	log.SetFlags(0)
	addr := make(chan string)
	go startServer(addr)
	// the client here can just refer to a computer such as my computer
	client, _ := geerpc.Dial("tcp", <-addr)
	defer func() { _ = client.Close() }()

	time.Sleep(time.Second)
	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			args := &Args{Num1: i, Num2: i * i}
			var reply int
			if err := client.Call("Foo.Sum", args, &reply); err != nil {
				log.Fatal("call Foo.Sum error:", err)
			}
			log.Printf("%d + %d = %d", args.Num1, args.Num2, reply)
		}(i)
	}

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			args := &Args{Num1: i, Num2: i * i}
			var reply int
			if err := client.Call("Foo.Minus", args, &reply); err != nil {
				log.Fatal("call Foo.Minus error:", err)
			}
			log.Printf("%d - %d = %d", args.Num2, args.Num1, reply)
		}(i)
	}

	wg.Wait()
}

// func startServer(addr chan string) {
// 	l, err := net.Listen("tcp", ":0")

// 	if err != nil {
// 		log.Fatal("network error:", err)
// 	}
// 	log.Println("start rpc server on", l.Addr())
// 	addr <- l.Addr().String()
// 	geerpc.Accept(l)
// }

// func main() {
//version day 2
// log.SetFlags(0)
// addr := make(chan string)
// go startServer(addr)

// client, _ := geerpc.Dial("tcp", <-addr) // this place can change the codectype
// defer func() { _ = client.Close() }()

// time.Sleep(time.Second)

// var wg sync.WaitGroup
// for i := 0; i < 5; i++ {
// 	wg.Add(1)
// 	go func(i int) {
// 		defer wg.Done()
// 		args := fmt.Sprintf("geerpc req %d", i)
// 		var reply string
// 		if err := client.Call("Foo.Sum", args, &reply); err != nil {
// 			log.Fatal("call Foo.Sum error:", err)
// 		}
// 		log.Println("reply:", reply)
// 	}(i)
// }
// wg.Wait()

// version day 1
// addr := make(chan string)
// go startServer(addr)

// conn, _ := net.Dial("tcp", <-addr)
// defer func() { _ = conn.Close() }()

// time.Sleep(time.Second)

// _ = json.NewEncoder(conn).Encode(geerpc.DefaultOption)

// cc := codec.NewGobCodec(conn)

// for i := 0; i < 5; i++ {
// 	h := &codec.Header{
// 		ServiceMethod: "Foo.Sum",
// 		Seq:           uint64(i),
// 	}

// 	_ = cc.Write(h, fmt.Sprintf("geerpc req %d", h.Seq))
// 	_ = cc.ReadHeader(h)
// 	var reply string
// 	_ = cc.ReadBody(&reply)
// 	log.Println("reply: ", reply)
// }

// }
