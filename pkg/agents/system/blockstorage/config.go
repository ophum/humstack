package blockstorage

type BlockStorageAgentDownloadAPIConfig struct {
	AdvertiseAddress string `yaml:"advertiseAddress"`
	ListenAddress    string `yaml:"listenAddress"`
	ListenPort       int32  `yaml:"listenPort"`
}
type BlockStorageAgentConfig struct {
	DownloadAPI         BlockStorageAgentDownloadAPIConfig `yaml:"downloadAPI"`
	BlockStorageDirPath string                             `yaml:"blockStorageDirPath"`
	ImageDirPath        string                             `yaml:"imageDirPath"`
	ParallelLimit       int64                              `yaml:"parallelLimit"`
}
