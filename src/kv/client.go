package main
//client
import(
	"net/rpc"
	"log"
	"fmt"
)

func connect() *rpc.Client {
	//Dial在指定的网络和地址和rpc服务器相连
	client, err := rpc.Dial("tcp", ":8080")
	if err != nil {
		//log.Fatal("client :", err)
		fmt.Println("error")
		return nil
	}
	return client
}

func get(key string) string {
	client := connect()
	args := GetArgs{key}
	reply := GetReply{}
	err := client.Call("KV.Get", args, &reply)
	if err != nil {
		log.Fatal(err)
	}
	client.Close()
	return reply.Value

}

func put(key string, val string){
	client := connect()
	args := PutArgs{key, val}
	reply := PutReply{}
	err := client.Call("KV.Put", args, &reply)
	if err != nil {
		log.Fatal(err)
	}
	client.Close()
}
