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
	"github.com/nuclio/nuclio/pkg/zap"
	"github.com/nuclio/nuclio-sdk"
	"github.com/v3io/v3io-go-http"
	"github.com/pkg/errors"
	"time"
	"fmt"
)

// Wrapper for invoking nuclio functions
//
// Usage example:
//
// func main()  {
// 	data  := nutest.DataBind{Name:"db0", Url:"<v3io-IP:Port>", Container:"<data-container-name>"}
// 	event := nutest.TestEvent{Body: []byte("test")}
// 	nutest.Invoke(MyHandler, nutest.TestSpec{
// 		Event:&event, Data:&data, LogLevel:nucliozap.InfoLevel})
// }
//
// func MyHandler(context *nuclio.Context, event nuclio.Event) (interface{}, error) {
// 	context.Logger.Debug("some text")
// 	return "resp", nil
// }
//

type TestSpec struct {
	LogLevel nucliozap.Level
	Event    *TestEvent
	Data     *DataBind
}

func Invoke(nfunc func(context *nuclio.Context, event nuclio.Event)(interface {}, error), spec TestSpec) error {
	//if spec.LogLevel == nil {
	//	spec.LogLevel = nucliozap.InfoLevel
	//}

	logger, err := nucliozap.NewNuclioZapCmd("emulator", spec.LogLevel)
	if err != nil {
		return errors.Wrap(err, "Failed to create logger")
	}

	db := map[string]nuclio.DataBinding{}
	if spec.Data != nil {
		container, err := createContainer(logger, spec.Data.Url, spec.Data.Container)
		if err != nil {
			logger.ErrorWith("Failed to createContainer", "err", err)
			return errors.Wrap(err, "Failed to createContainer")
		}

		if spec.Data.Name == "" {
			spec.Data.Name = "db0"
		}
		db[spec.Data.Name] = container
	}

	context := nuclio.Context{Logger:logger, DataBinding:db}

	body, err := nfunc(&context, spec.Event)
	if err != nil {
		logger.ErrorWith("Function execution failed", "err", err)
		return err
	}
	logger.InfoWith("Function completed","output",body)
	fmt.Println(body)

	return nil
}

func createContainer(logger nuclio.Logger, addr, cont string) (*v3io.Container, error) {
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

type TestEvent struct {
	Body               []byte
	ContentType        string
	sourceInfoProvider nuclio.SourceInfoProvider
	id                 nuclio.ID
	emptyByteArray     []byte
	emptyHeaders       map[string]interface{}
	emptyTime          time.Time
}

var ErrUnsupported = errors.New("Event does not support this interface")

func (te *TestEvent) GetVersion() int {
	return 0
}

func (te *TestEvent) SetSourceProvider(sourceInfoProvider nuclio.SourceInfoProvider) {
	te.sourceInfoProvider = sourceInfoProvider
}

func (te *TestEvent) GetSource() nuclio.SourceInfoProvider {
	return te.sourceInfoProvider
}

func (te *TestEvent) GetID() nuclio.ID {
	return te.id
}

func (te *TestEvent) SetID(id nuclio.ID) {
	te.id = id
}

func (te *TestEvent) GetContentType() string {
	return te.ContentType
}

func (te *TestEvent) GetBody() []byte {
	return te.Body
}

func (te *TestEvent) GetSize() int {
	return 0
}

func (te *TestEvent) GetHeader(key string) interface{} {
	return nil
}

func (te *TestEvent) GetHeaderByteSlice(key string) []byte {
	return te.emptyByteArray
}

func (te *TestEvent) GetHeaderString(key string) string {
	return string(te.GetHeaderByteSlice(key))
}

func (te *TestEvent) GetHeaders() map[string]interface{} {
	return te.emptyHeaders
}

func (te *TestEvent) GetTimestamp() time.Time {
	return te.emptyTime
}

func (te *TestEvent) GetPath() string {
	return ""
}

func (te *TestEvent) GetURL() string {
	return ""
}

func (te *TestEvent) GetMethod() string {
	return ""
}

func (te *TestEvent) GetField(key string) interface{} {
	return nil
}

func (te *TestEvent) GetFieldByteSlice(key string) []byte {
	return nil
}

func (te *TestEvent) GetFieldString(key string) string {
	return ""
}

func (te *TestEvent) GetFieldInt(key string) (int, error) {
	return 0, ErrUnsupported
}

func (te *TestEvent) GetFields() map[string]interface{} {
	return nil
}



