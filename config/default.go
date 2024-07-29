package config

// DefaultValues is the default configuration
const DefaultValues = `
[Log]
Environment = "development" # "production" or "development"
Level = "info"
Outputs = ["stderr"]

[Etherman]
URL = "http://localhost:8545"

[EthTxManager]
FrequencyToMonitorTxs = "1s"
WaitTxToBeMined = "2m"
ForcedGas = 0
GasPriceMarginFactor = 1
MaxGasPriceLimit = 0

[SequenceSender]
WaitPeriodSendSequence = "5s"
MaxTxSizeForL1 = 131072
L2Coinbase = "0xf39fd6e51aad88f6f4ce6ab8827279cfffb92266"
GasOffset = 80000
MaxBatchesForL1 = 10
DAPermitApiPrivateKey = {Path = "/pk/sequencer.keystore", Password = "testonly"}

[DataAvailability]
NubitRpcURL = "http://127.0.0.1:26658"
NubitAuthKey = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJBbGxvdyI6WyJwdWJsaWMiLCJyZWFkIiwid3JpdGUiLCJhZG1pbiJdfQ.DAMv0s7915Ahx-kDFSzDT1ATz4Q9WwktWcHmjp7_99Q"
NubitNamespace = "xlayer"
NubitMaxBatchesSize = "102400"
NubitGetProofMaxRetry = "10"
NubitGetProofWaitPeriod = "5s"
`
