package nodenetwork

type NetworkAgentConfigVXLAN struct {
	DevName string `yaml:"devName"`
	Group   string `yaml:"group"`
}

type NetworkAgentConfigVLAN struct {
	DevName                 string `yaml:"devName"`
	VLANInterfaceNamePrefix string `yaml:"vlanInterfaceNamePrefix"`
}

type NetworkAgentConfig struct {
	VXLAN NetworkAgentConfigVXLAN `yaml:"vxlan"`
	VLAN  NetworkAgentConfigVLAN  `yaml:"vlan"`
}
