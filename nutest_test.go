package nutest

import (
	"testing"
	"github.com/nuclio/nuclio-sdk-go"
	"fmt"
)

func MyHandler(context *nuclio.Context, event nuclio.Event) (interface{}, error) {
	context.Logger.Info("path: %s", string(event.GetPath()))
	context.Logger.Info("body: %s", string(event.GetBody()))
	context.Logger.Info("headers: %+v", event.GetHeaders())
	context.Logger.Info("str header: %s", event.GetHeaderString("first"))
	return "test me\n" + string(event.GetBody()), nil
}


func TestName(t *testing.T) {
	data := DataBind{Name:"db0", Url:"<v3io address>", Container:"x"}
	tc, err := NewTestContext(MyHandler, true, &data )
	if err != nil {
		t.Fail()
	}

	testEvent := TestEvent{
		Path: "/some/path",
		Body: []byte("1234"),
		Headers:map[string]interface{}{"first": "string", "sec": "1"},
		}
	resp, err := tc.Invoke(&testEvent)
	fmt.Println("resp:", resp)
	fmt.Println("err:", err)
}

