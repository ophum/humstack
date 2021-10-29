# humstack v1

humstack is iaas. influenced by n0stack, kubernetes...

**IN PROGRESS...**

## development
### migrate db

```
$ go install github.com/rubenv/sql-migrate/...@master
$ sql-migrate up
```
### run api server
```
$ make run-api
```

### run disk agent
```
$ make run-disk-agent
```

### humcli

#### list disks
```
$ go run cmd/humcli/main.go list disks
name    size    status
hoge    1G      Active
hoge2   1G      Active
```
#### create disk
```
$ go run cmd/humcli/main.go create disk --name test --limit 10G
name    size    status
test    10G      Pending
```