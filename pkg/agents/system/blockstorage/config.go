package blockstorage

type BlockStorageAgentDownloadAPIConfig struct {
	AdvertiseAddress string `yaml:"advertiseAddress"`
	ListenAddress    string `yaml:"listenAddress"`
	ListenPort       int32  `yaml:"listenPort"`
}

type BlockStorageAgentCephBackendConfig struct {
	ConfigPath string `yaml:"configPath"`
	PoolName   string `yaml:"poolName"`
}
type BlockStorageAgentConfig struct {
	DownloadAPI         BlockStorageAgentDownloadAPIConfig  `yaml:"downloadAPI"`
	BlockStorageDirPath string                              `yaml:"blockStorageDirPath"`
	ImageDirPath        string                              `yaml:"imageDirPath"`
	CephBackend         *BlockStorageAgentCephBackendConfig `yaml:"cephBackend"`
	ParallelLimit       int64                               `yaml:"parallelLimit"`
}
