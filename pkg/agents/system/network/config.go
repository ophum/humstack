package network

type NetworkAgentConfigVXLAN struct {
	DevName string `yaml:"devName"`
	Group   string `yaml:"group"`
}

type NetworkAgentConfig struct {
	VXLAN NetworkAgentConfigVXLAN `yaml:"vxlan"`
}
