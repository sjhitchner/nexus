package filter

import (
	. "github.com/sjhitchner/nexus/common"
	"testing"
	"time"
)

var js = `
var buckets = ["mybucket", "mybucket2"];

var filter = function (bucket, value, raw) {
	console.log(value.Type);
	Emit(bucket, value, raw);
}
`

func TestFilter(t *testing.T) {

	filter, _ := NewJSFilter("testfilter", js)

	payload := Payload{
		Type:      "test",
		Data:      "data",
		CreatedAt: time.Now(),
	}

	filter.Filter(payload)
	filter.Filter(payload)
	filter.Filter(payload)
}
