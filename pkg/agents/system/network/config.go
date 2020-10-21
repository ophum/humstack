package network

type NetworkAgentConfigVXLAN struct {
	DevName string `yaml:"devName"`
	Group   string `yaml:"group"`
}

type NetworkAgentConfigVLAN struct {
	DevName string `yaml:"devName"`
}

type NetworkAgentConfig struct {
	VXLAN NetworkAgentConfigVXLAN `yaml:"vxlan"`
	VLAN  NetworkAgentConfigVLAN  `yaml:"vlan"`
}
