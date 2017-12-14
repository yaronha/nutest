# nuclio function wrapper

Test nuclio functions locally 

# Usage:

```golang
package main

import (
	"github.com/yaronha/nutest"
	"github.com/nuclio/nuclio-sdk"
	"github.com/nuclio/nuclio/pkg/zap"
)


func main()  {
	// data binding for V3IO data containers, optional 
	data := nutest.DataBind{Name:"db0", Url:Url:"<v3io-IP:Port>", Container:"<data-container-name>"}
	// event data 
	event := nutest.TestEvent{Body: []byte("test")}
	
	nutest.Invoke(MyHandler, nutest.TestSpec{
		Event:&event, Data:&data, LogLevel:nucliozap.InfoLevel})
}

// nuclio function 
func MyHandler(context *nuclio.Context, event nuclio.Event) (interface{}, error) {
	context.Logger.Debug("test")
	return "resp", nil
}

```