package ports

import "gitag.ir/cookthepot/services/vault/modules/logger/engines"

type LoggerEngineType interface {
	*engines.LoggerEngineFile | *engines.LoggerEngineStdout | any
}
