# humstack

humstack is iaas. influenced by n0stack, kubernetes...

## setup

依存パッケージ

```
sudo apt update
sudo apt install qemu qemu-kvm cloud-image-utils
echo 1 > /proc/sys/ipv4/ip_forward
```

## ビルド

```
make all
```

## 実行

### apiserver

```
./apiserver --listen-address 0.0.0.0 --listen-port 8080
```

### agent

管理者権限で実行する。実行したマシンのホスト名が node 名として apiserver に登録される。

```
sudo ./agent --config config.yaml
```

#### config.yaml

```
# apiserverのアドレス
apiServerAddress: localhost
# apiserverのポート
apiServerPort: 8080

# agentが動作するノードのリソース量(使われていない)
limitMemory: 8G
limitVcpus: 8000m

# blockstorageを保存する場所
blockStorageDirPath: ./blockstorages

# networkAgentの設定
networkAgentConfig:
  # vxlanの設定
  vxlan:
    # デバイス名
    devName: eth0
    # vxlanで使用するマルチキャストIP
    group: 239.0.0.1

```

## humcli

yaml ファイルを読み込んで apiserver にリクエストを送信するコマンドラインツール

```
humstack cli

Usage:
  humstack [command]

Available Commands:
  create
  delete
  get
  help        Help about any command
  update
  watch

Flags:
      --api-server-address string   apiserver address (default "localhost")
      --api-server-port int32       apiserver Port (default 8080)
      --config string               config file
      --g string                    group id (default "default")
  -h, --help                        help for humstack
      --n string                    namespace id (default "default")

Use "humstack [command] --help" for more information about a command.
```

### リソース

#### corev0/group

グループ、組織

```
meta:
  apiType: corev0/group
  id: group1
  name: group1
```

#### corev0/namespace

グループ内でリソースを分離

```
meta:
  apiType: corev0/namespace
  id: ns1
  name: namespace1
  group: group1
```

#### corev0/externalippool

外部ネットワークの設定。group や namespace は指定しない

```
meta:
  apiType: corev0/externalippool
  id: eippool
  name: eippool
spec:
  ipv4CIDR: 192.168.10.0/24
  BridgeName: exBr
  DefaultGateway: 192.168.10.254
```

#### corev0/externalip

外部ネットワークのアドレス。

```
meta:
  apiType: corev0/externalip
  id: eip1
  name: eip1
  group: group1
  namespace: ns1
spec:
  poolID: eippool
  ipv4Address: 192.168.10.100
  ipv4Prefix: 24
```

#### systemv0/network

仮想ネットワーク。Linux Bridge や vxlan などが作成される。

```
meta:
  apiType: systemv0/network
  id: net1
  name: network1
  group: group1
  namespace: ns1
  annotations:
    networkv0/network_type: VXLAN
spec:
  # vxlanやvlanで使用するID
  id: "100"
  # そのネットワークのCIDR
  ipv4CIDR: 10.0.0.0/24
```

##### annotations

| key                       | value                     | description                                                                                          |
| ------------------------- | ------------------------- | ---------------------------------------------------------------------------------------------------- |
| networkv0/network_type    | `VXLAN`, `VLAN`, `Bridge` | `VXLAN`の場合`vxlan`の link と Bridge が作成される。`VXLAN` 未実装。`Bridge`は Bridge のみ作成される |
| networkv0/bridge_name     |                           | agent によって作成された Bridge の名前が入る                                                         |
| networkv0/default_gateway | `xxx.xxx.xxx.xxx/xx`      | 指定されたアドレスが Bridge に対して設定され、コンピュートノード上の iptables で NAPT される         |

#### systemv0/virtualrouter

仮想ルーター。指定したノード上で netns と iptables などを利用したルーティング、NAT を行う。

```
meta:
  apiType: systemv0/virtualrouter
  id: vrouter1
  name: virtualrouter1
  group: group1
  namespace: ns1
  annotations:
    virtualrouterv0/node_name: worker2
spec:
  externalGateway: 192.168.10.254
  externalIPs:
    - externalIPID: eip1
      bindInternalIPv4Address: 10.0.0.1
  natGatewayIP: 192.168.10.200
  nics:
    - networkID: net1
      ipv4Address: 10.0.0.254/24

```

#### systemv0/blockstorage

仮想ディスク。

```
meta:
  apiType: systemv0/blockstorage
  id: bs1
  name: blockstorage1
  group: group1
  namespace: ns1
  annotations:
    blockstoragev0/node_name: worker1
    blockstoragev0/type: Local
spec:
  requestSize: 1G
  limitSize: 10G
  from:
    type: HTTP
    http:
      url: http://192.168.20.2:8082/focal-server-cloudimg-amd64.img

```

#### systemv0/virtualmachine

仮想マシン。blockstorage や network などに依存するため、それらが利用できる状態になるまで作成されない。

```
meta:
  apiType: systemv0/virtualmachine
  id: vm1
  name: virtualmachine1
  group: group1
  namespace: ns1
  annotations:
    virtualmachinev0/node_name: worker1
spec:
  requestVcpus: 1000m
  limitVcpus: 1000m
  requestMemory: 1G
  limitMemory: 1G
  blockStorageIDs:
    - bs1
  nics:
    - networkID: net1
      ipv4Address: 10.0.0.1
      nameservers:
        - 8.8.8.8
      defaultGateway: 10.0.0.254
  actionState: PowerOn
  loginUsers:
    - username: test
      sshAuthorizedKeys:
        - ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQCsf7CDppU1lSzUbsmszAXX/rAXdGxB71i93IsZtV4omO/uRz/z6dLIsBidf9vIqcEfCFTFR00ULC+GKULTNz2LOaGnGsDS28Bi5u+cx90+BCAzEg6cBwPIYmdZgASsjMmRvI/r+xR/gNxq2RCR8Gl8y5voAWoU8aezRUxf1Ra3KljMd1dbIFGJxgzNiwqN3yL0tr9zActw/Q7yBWKWi1c5sW2QZLAnSj/WWTSGGm0Ad88Aq22DakwN6itUkS6XNhr4YKehLVm90fIojrCrtZmClULAlnUk5lbdzou4jiETsZz3zk/q76ZQ3ugk+G00kcx9v6ElLkAFv2ZZqzWbMvUz6J0k2SzkAIbcBDz+aq2sXeY04FaIOFPiH41+DTQXCtOskWkaJBMKLTE/Z83nSyQGr9If2F/PbnuxGkwiZzeZaLWxqI2SebhLR5jPETgfhB1y83RP6u8Jq5+9BUURFqpb8mfG/riTnAj0ZR4Li23+/hWhc8We+fVB1BxdbWyRn/M=
```
