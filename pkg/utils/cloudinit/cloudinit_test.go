package cloudinit

import "testing"

func TestOutput(t *testing.T) {
	metaData := MetaData{
		InstanceID:    "test",
		LocalHostName: "test",
	}
	userData := UserData{
		Users: []UserDataUser{
			{
				Name:   "testuser",
				Groups: "sudo",
				Shell:  "/bin/bash",
				Sudo: []string{
					"ALL=(ALL) NOPASSWD:ALL",
				},
				LockPasswd: true,
			},
		},
	}
	networkConfig := NetworkConfig{
		Version: 1,
		Config: []NetworkConfigConfig{
			{
				Type:       NetworkConfigConfigTypePhysical,
				Name:       "eth0",
				MacAddress: "52:54:00:01:02:03",
			},
		},
	}

	ci := NewCloudInit(metaData, userData, networkConfig)
	err := ci.Output("./")
	if err != nil {
		t.Fatal(err)
	}
}
