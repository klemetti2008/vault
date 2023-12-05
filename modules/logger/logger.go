package logger

import (
	"reflect"

	"gitag.ir/cookthepot/services/vault/modules/logger/engines"
	"gitag.ir/cookthepot/services/vault/modules/logger/ports"
)

type Logger[T ports.LoggerEngineType] struct {
	Engine LoggerEngine[T]
}

var Enable = false
var Default *Logger[ports.LoggerEngineType] = Initialize[ports.LoggerEngineType]()

func Initialize[T ports.LoggerEngineType](engines ...T) *Logger[T] {
	var loggerTypes []string = []string{}
	for _, engine := range engines {
		loggerType := reflect.ValueOf(&engine).Elem().Type().Name()
		loggerTypes = append(loggerTypes, loggerType)
	}
	var Default = &Logger[T]{}
	Default.Engine.Instances = &engines
	Default.Engine.Types = loggerTypes
	return Default
}
func (l Logger[T]) GetInstance() *Logger[ports.LoggerEngineType] {
	return Default
}
func (l Logger[T]) Trace(args ...interface{}) {
	log("Trace", args...)
}
func (l Logger[T]) Debug(args ...interface{}) {
	log("Debug", args...)
}
func (l Logger[T]) Info(args ...interface{}) {
	log("Info", args...)
}
func (l Logger[T]) Warn(args ...interface{}) {
	log("Warn", args...)
}
func (l Logger[T]) Error(args ...interface{}) {
	log("Error", args...)
}
func (l Logger[T]) Fatal(args ...interface{}) {
	log("Fatal", args...)
}
func log(functionName string, args ...interface{}) {
	if !Enable {
		return
	}
	for _, instance := range *Default.Engine.Instances {
		loggerFile, ok := instance.(*engines.LoggerEngineFile)
		if ok {
			method := getFunction(loggerFile, functionName)
			go method(args...)
		}
		loggerStdout, ok := instance.(*engines.LoggerEngineStdout)
		if ok {
			method := getFunction(loggerStdout, functionName)
			go method(args...)
		}
	}
}

func getFunction[T ports.LoggerEngineType](obj T, functionName string) func(...interface{}) {
	methodVal := reflect.ValueOf(obj).MethodByName(functionName)
	methodInterface := methodVal.Interface()
	method := methodInterface.(func(...interface{}))
	return method
}
