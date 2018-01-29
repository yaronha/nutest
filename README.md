# nuclio function wrapper

Test nuclio functions locally 

# Usage:

```golang
package main

import (
	"github.com/nuclio/nuclio-sdk-go"
	"fmt"
)

func main() {
	// data binding for V3IO data containers, optional 
	data := DataBind{Name:"db0", Url:"192.168.1.1", Container:"x"}

	// Create TestContext and specify the function name, verbose, data 
	tc, err := NewTestContext(MyHandler, true, &data )
	if err != nil {
		panic(err)
	}

	// Create a new test event 
	testEvent := TestEvent{
		Path: "/some/path",
		Body: []byte("1234"),
		Headers:map[string]interface{}{"first": "something"},
	}
	
	// invoke the tested function with the new event and print it's output 
	resp, err := tc.Invoke(&testEvent)
	fmt.Println("resp:", resp)
	fmt.Println("err:", err)
}
```