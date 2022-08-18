启动etcd默认服务器 https://github.com/etcd-io/etcd/releases
启动nats默认服务器 https://github.com/nats-io/nats-server/releases
go run main.go -port=3250
go run main.go -port=3251 -type=room -frontend=false

cli命令
go run ./cli

