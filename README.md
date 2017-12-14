# nuclio function wrapper

Test nuclio functions locally 

# Usage:

```golang
package main

import (
	"github.com/yaronha/nutest"
	"github.com/nuclio/nuclio-sdk"
)


func main()  {
	// data binding for V3IO data containers, optional, can use nil instead 
	data := nutest.DataBind{Name:"db0", Url:Url:"<v3io-IP:Port>", Container:"<data-container-name>"}
	// event data 
	event := nutest.TestEvent{Body: []byte("test")}
	
	nutest.Invoke(MyHandler, &event, &data)
}

// nuclio function 
func MyHandler(context *nuclio.Context, event nuclio.Event) (interface{}, error) {
	context.Logger.Debug("test")
	return "resp", nil
}

```