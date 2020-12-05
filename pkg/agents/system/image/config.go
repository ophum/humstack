package image

type ImageAgentDownloadAPIConfig struct {
	AdvertiseAddress string `yaml:"advertiseAddress"`
	ListenPort       int32  `yaml:"listenPort"`
}

type ImageAgentCephBackendConfig struct {
	ConfigPath string `yaml:"configPath"`
	PoolName   string `yaml:"poolName"`
}
type ImageAgentConfig struct {
	ImageDirPath        string                      `yaml:"imageDirPath"`
	BlockStorageDirPath string                      `yaml:"blockStorageDirPath"`
	DownloadAPI         ImageAgentDownloadAPIConfig `yaml:"downloadAPI"`
	CephBackend         ImageAgentCephBackendConfig `yaml:"cephBackend"`
}
