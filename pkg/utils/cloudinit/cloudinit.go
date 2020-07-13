package cloudinit

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

type MetaData struct {
	InstanceID    string `yaml:"instance-id"`
	LocalHostName string `yaml:"local-hostname"`
}

type UserDataUser struct {
	Name              string   `yaml:"name"`
	Groups            string   `yaml:"groups"`
	Shell             string   `yaml:"shell"`
	Sudo              []string `yaml:"sudo"`
	SSHAuthorizedKeys []string `yaml:"ssh-authorized-keys"`
	LockPasswd        bool     `yaml:"lock_passwd"`
}

type UserData struct {
	Users []UserDataUser `yaml:"users"`
}

type NetworkConfigConfigSubnetType string

const (
	NetworkConfigConfigSubnetTypeStatic NetworkConfigConfigSubnetType = "static"
)

type NetworkConfigConfigSubnet struct {
	Type        NetworkConfigConfigSubnetType `yaml:"type"`
	Address     string                        `yaml:"address"`
	Netmask     string                        `yaml:"netmask"`
	Gateway     string                        `yaml:"gateway"`
	Nameservers []string                      `yaml:"dns_nameservers"`
}

type NetworkConfigConfigType string

const (
	NetworkConfigConfigTypePhysical NetworkConfigConfigType = "physical"
)

type NetworkConfigConfig struct {
	Type       NetworkConfigConfigType     `yaml:"type"`
	Name       string                      `yaml:"name"`
	MacAddress string                      `yaml:"mac_address"`
	Subnets    []NetworkConfigConfigSubnet `yaml:"subnets"`
}

type NetworkConfig struct {
	Version int32                 `yaml:"version"`
	Config  []NetworkConfigConfig `yaml:"config"`
}

type CloudInit struct {
	metaData      MetaData
	userData      UserData
	networkConfig NetworkConfig
}

func NewCloudInit(metaData MetaData, userData UserData, networkConfig NetworkConfig) *CloudInit {
	return &CloudInit{
		metaData:      metaData,
		userData:      userData,
		networkConfig: networkConfig,
	}
}

func (c *CloudInit) Output(dirPath string) error {
	name := c.metaData.InstanceID
	path := filepath.Join(dirPath, name)
	if !fileIsExists(path) {
		err := os.MkdirAll(path, 0755)
		if err != nil {
			return err
		}
	}

	metaDataYAML, err := yaml.Marshal(c.metaData)
	if err != nil {
		return err
	}

	userDataYAML, err := yaml.Marshal(c.userData)
	if err != nil {
		return err
	}

	networkConfigYAML, err := yaml.Marshal(c.networkConfig)
	if err != nil {
		return err
	}

	metaDataPath := filepath.Join(path, "meta-data")
	err = ioutil.WriteFile(metaDataPath, metaDataYAML, 0666)
	if err != nil {
		return err
	}

	userDataPath := filepath.Join(path, "user-data")
	userDataYAML = []byte(fmt.Sprintf("#cloud-config\n%s", userDataYAML))
	err = ioutil.WriteFile(userDataPath, userDataYAML, 0666)
	if err != nil {
		return err
	}

	networkConfigPath := filepath.Join(path, "network-config")
	err = ioutil.WriteFile(networkConfigPath, networkConfigYAML, 0666)
	if err != nil {
		return err
	}

	command := "cloud-localds"
	args := []string{
		"-N",
		networkConfigPath,
		filepath.Join(path, "cloudinit.img"),
		userDataPath,
		metaDataPath,
	}

	cmd := exec.Command(command, args...)
	if _, err := cmd.CombinedOutput(); err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

func fileIsExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
