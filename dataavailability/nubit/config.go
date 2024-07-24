package nubit

import "time"

type Config struct {
	NubitRpcURL             string        `mapstructure:"NubitRpcURL"`
	NubitModularAppName     string        `mapstructure:"NubitModularAppName"`
	NubitAuthKey            string        `mapstructure:"NubitAuthKey"`
	NubitNamespace          string        `mapstructure:"NubitNamespace"`
	NubitMaxBatchesSize     uint64        `mapstructure:"NubitMaxBatchesSize"`
	NubitGetProofMaxRetry   uint64        `mapstructure:"NubitGetProofMaxRetry"`
	NubitGetProofWaitPeriod time.Duration `mapstructure:"NubitGetProofWaitPeriod"`
}
