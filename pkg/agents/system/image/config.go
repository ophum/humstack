package image

type ImageAgentDownloadAPIConfig struct {
	AdvertiseAddress string `yaml:"advertiseAddress"`
	ListenPort       int32  `yaml:"listenPort"`
}
type ImageAgentConfig struct {
	ImageDirPath        string                      `yaml:"imageDirPath"`
	BlockStorageDirPath string                      `yaml:"blockStorageDirPath"`
	DownloadAPI         ImageAgentDownloadAPIConfig `yaml:"downloadAPI"`
}
