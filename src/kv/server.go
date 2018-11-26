package main 
import(
	"log"
	"net/rpc"
	"sync"
	"net"
	"fmt"
)

type KV struct {
	mu sync.Mutex
	data map[string]string
}

func server(){
	kv := new(KV)
	kv.data = make(map[string]string)
	rpcs := rpc.NewServer()
	rpcs.Register(kv)
	l, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}
	go func(){
		for {
			conn, err := l.Accept()
			if err == nil {
				fmt.Printf("recieve a req\n")
				go rpcs.ServeConn(conn)
			}else {
				break
			}
		}
		l.Close()
	}()
	fmt.Println("server start at 8080")
}

func (kv *KV) Get(args GetArgs, reply *GetReply) error{
	kv.mu.Lock()
	defer kv.mu.Unlock()

	val, ok := kv.data[args.Key]
	if ok {
		reply.Err = OK
		reply.Value = val
	}else {
		reply.Err = ErrNoKey
		reply.Value = ""
	}
	return nil
}

func (kv *KV) Put(args PutArgs, reply *PutReply)error{
	kv.mu.Lock()
	defer kv.mu.Unlock()

	kv.data[args.Key] = args.Value
	fmt.Printf("put %s -> %s \n", args.Key, args.Value)
	reply.Err = OK
	return nil
}
