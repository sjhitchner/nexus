package filter

import (
	"encoding/json"
	"fmt"
	"github.com/robertkrimen/otto"
	. "github.com/sjhitchner/nexus/interfaces"
	"github.com/sjhitchner/nexus/interfaces/router"
)

type JSFilter struct {
	name      string
	buckets   []Bucket
	vm        *otto.Otto
	isEnabled bool
}

func NewJSFilter(name string, js string) (Filter, error) {
	filter := &JSFilter{
		name:      name,
		isEnabled: true,
		vm:        otto.New(),
	}
	filter.vm.Set("Emit", filter.Emit)

	script, err := filter.vm.Compile("", js)
	if err != nil {
		return nil, err
	}

	_, err = filter.vm.Run(script)
	if err != nil {
		return nil, err
	}

	filter.buckets, err = getBucketList(filter.vm)
	if err != nil {
		return nil, err
	}

	return filter, nil
}

func (t *JSFilter) Filter(payload Payload) error {
	b, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	payloadObj, err := t.vm.ToValue(payload)
	if err != nil {
		return err
	}

	_, err = t.vm.Call("filter", nil, "mybucket", payloadObj, string(b))
	if err != nil {
		return err
	}

	return nil
}

func (t *JSFilter) Enable() {
	t.isEnabled = true
}

func (t *JSFilter) Disable() {
	t.isEnabled = false
}

func (t *JSFilter) IsEnabled() bool {
	return t.isEnabled
}

func (t *JSFilter) Name() string {
	return t.name
}

func (t *JSFilter) Buckets() []Bucket {
	return t.buckets
}

func (t *JSFilter) Emit(call otto.FunctionCall) otto.Value {
	path := call.Argument(0).String()
	data, err := call.Argument(1).Export()
	if err == nil {
		router.Route(path, data)
	}
	return otto.Value{}
}

func getBucketList(vm *otto.Otto) ([]Bucket, error) {
	bucketValue, err := vm.Get("buckets")
	if err != nil {
		return nil, err
	}

	bucketsInter, err := bucketValue.Export()
	if err != nil {
		return nil, err
	}

	switch val := bucketsInter.(type) {
	case string:
		return []Bucket{Bucket(val)}, nil
	case []interface{}:
		arr := make([]Bucket, len(val))
		for i, inter := range val {
			str, ok := inter.(string)
			if !ok {
				return nil, fmt.Errorf("invalid bucket name")
			}
			arr[i] = Bucket(str)
		}
		return arr, nil
	default:
		return nil, fmt.Errorf("invalid bucket list")
	}
}
