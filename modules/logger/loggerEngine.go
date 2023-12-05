package logger

import "gitag.ir/cookthepot/services/vault/modules/logger/ports"

type LoggerEngine[T ports.LoggerEngineType] struct {
	Types     []string
	Instances *[]T
}
