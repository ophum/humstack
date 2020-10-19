package blockstorage

type BlockStorageAgentDownloadAPIConfig struct {
	ListenAddress string `yaml:"listenAddress"`
	ListenPort    int32  `yaml:"listenPort"`
}
type BlockStorageAgentConfig struct {
	DownloadAPI         BlockStorageAgentDownloadAPIConfig `yaml:"downloadAPI"`
	BlockStorageDirPath string                             `yaml:"blockStorageDirPath"`
}
