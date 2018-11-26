package main
import(
	"fmt"
	"testing"
)

func TestKV(t *testing.T){
	server()
	put("yihau", "go")
	fmt.Println(get("yihau"))
}
