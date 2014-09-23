package main

import (
	//"encoding/json"
	"fmt"
	//"github.com/robertkrimen/otto"
	//. "github.com/sjhitchner/infosphere/common"
	"github.com/sjhitchner/infosphere/interfaces/router"
	"time"
	//"io/ioutil"
	//"code.google.com/p/go.net/websocket"
	"math/rand"
)

func main() {

	path := "test"

	for i := 0; i < 20; i++ {
		go worker(i, path)
	}

	for {
		select {
		case <-time.After(1 * time.Millisecond):
			for i := 0; i < 1000; i++ {
				router.Route(path, i)
			}
		case <-time.After(30 * time.Second):
			return
		}
	}
}

func worker(id int, path string) {
	channel := router.NewChannel(10)
	router.AddChannel(path, channel)

	go func() {
		dur := time.Duration(rand.Intn(30))
		select {
		case <-time.After(dur * time.Second):
			channel.Close()
		}
	}()

	for value := range channel.Dequeue() {
		fmt.Println(value)
	}

	fmt.Printf("go routine ending [%d]\n", id)
}

/*


func main() {
	fmt.Println("Weeeeee")


		jsBytes, err := ioutil.ReadFile("filter.js")
		if err != nil {
			panic(err)
		}

	js := `
var buckets = ["mybucket", "mybucket2"];

var filter = function (bucket, value, raw) {
	console.log(value.Type);
	Emit(bucket, value, raw);
}
`

	vm := otto.New()
	vm.Set("Emit", Emit)

	script, err := vm.Compile("", js)
	if err != nil {
		fmt.Println(err)
		return
	}

	_, err = vm.Run(script)
	if err != nil {
		panic(err)
	}

	buckets, err := GetBuckets(vm)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Buckets: %v\n", buckets)

	payload := Payload{
		Type:      "test",
		Data:      "data",
		CreatedAt: time.Now(),
	}
	b, err := json.Marshal(payload)
	if err != nil {
		panic(err)
	}

	payloadObj, err := vm.ToValue(payload)
	if err != nil {
		panic(err)
	}

	_, err = vm.Call("filter", nil, "mybucket", payloadObj, string(b))
	if err != nil {
		panic(err)
	}

}

func Emit(call otto.FunctionCall) otto.Value {
	fmt.Println("emitting...",
		call.Argument(0).String(),
		call.Argument(1).String(),
		call.Argument(2).String())
	return otto.Value{}
}

func GetBuckets(vm *otto.Otto) ([]string, error) {
	bucketValue, err := vm.Get("buckets")
	if err != nil {
		return nil, err
	}

	fmt.Println("GB", bucketValue)

	bucketsInter, err := bucketValue.Export()
	if err != nil {
		return nil, err
	}

	fmt.Println("GBI", bucketsInter)

	switch val := bucketsInter.(type) {
	case string:
		fmt.Println("String", val)
		return []string{val}, nil
	case []interface{}:
		arr := make([]string, len(val))
		for i, inter := range val {
			str, ok := inter.(string)
			if !ok {
				return nil, fmt.Errorf("invalid bucket name")
			}
			arr[i] = str
		}
		return arr, nil
	default:
		return nil, fmt.Errorf("invalid bucket list")
	}
}
*/
