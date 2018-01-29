/*
Copyright 2017 The Nuclio Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package nutest

import (
	"github.com/nuclio/zap"
	"github.com/nuclio/nuclio-sdk-go"
	"github.com/v3io/v3io-go-http"
	"github.com/pkg/errors"
	"github.com/nuclio/logger"
)

/* Wrapper for invoking nuclio functions

Usage example:

func main() {
	data := DataBind{Name:"db0", Url:"192.168.1.1", Container:"x"}
	tc, err := NewTestContext(MyHandler, true, &data )
	if err != nil {
		panic(err)
	}

	testEvent := TestEvent{
		Path: "/some/path",
		Body: []byte("1234"),
		Headers:map[string]interface{}{"first": "something"},
	}
	resp, err := tc.Invoke(&testEvent)
	fmt.Println("resp:", resp)
	fmt.Println("err:", err)
}

*/

func NewTestContext(function func(context *nuclio.Context, event nuclio.Event)(interface {}, error),
	   verbose bool, data  *DataBind) (*TestContext, error) {
	newTest := TestContext{Data:data}
	if verbose {
		newTest.LogLevel = nucliozap.DebugLevel
	} else {
		newTest.LogLevel = nucliozap.InfoLevel
	}

	logger, err := nucliozap.NewNuclioZapCmd("emulator", newTest.LogLevel)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create logger")
	}

	newTest.logger = logger

	db := map[string]nuclio.DataBinding{}
	if data != nil {
		container, err := createContainer(logger, data.Url, data.Container)
		if err != nil {
			logger.ErrorWith("Failed to createContainer", "err", err)
			return nil, errors.Wrap(err, "Failed to createContainer")
		}

		if data.Name == "" {
			data.Name = "db0"
		}
		db[data.Name] = container
	}

	newTest.context = nuclio.Context{Logger:logger, DataBinding:db}
	newTest.function = function


	return &newTest, nil
}

type TestContext struct {
	LogLevel  nucliozap.Level
	logger    logger.Logger
	Data      *DataBind
	context   nuclio.Context
	function  func(context *nuclio.Context, event nuclio.Event)(interface {}, error)
}

func (tc *TestContext) Invoke(event nuclio.Event) (interface{}, error) {
	//if spec.LogLevel == nil {
	//	spec.LogLevel = nucliozap.InfoLevel
	//}


	body, err := tc.function(&tc.context, event)
	if err != nil {
		tc.logger.ErrorWith("Function execution failed", "err", err)
		return body, err
	}
	tc.logger.InfoWith("Function completed","output",body)

	return body, err
}

func createContainer(logger logger.Logger, addr, cont string) (*v3io.Container, error) {
	// create context
	context, err := v3io.NewContext(logger, addr , 8)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create client")
	}

	// create session
	session, err := context.NewSession("", "", "v3test")
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create session")
	}

	// create the container
	container, err := session.NewContainer(cont)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create container")
	}

	return container, nil
}

type DataBind struct {
	Name        string
	Url         string
	Container   string
	User        string
	Password    string
}
