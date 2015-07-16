# zkcli

a simple zookeeper client

# build

go build zkcli.go

# usage

    ./zkcli -h
    
    Usage:
     ./zkcli [OPTIONS] CMD PATH [CmdArgs]
    
    OPTIONS:
     --servers: required, zk servers
     --chroot: optional, zk chroot
     -r: optional, recursive operation ?
    
    CMD: [ls|children|get|set|create|delete|exists]
    
    PATH: zookeeper path
