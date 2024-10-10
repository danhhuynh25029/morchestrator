```
.
├── cmd
│   └── main.go
├── go.mod
├── go.sum
└── internal
    ├── manager
    │   └── manager.go
    ├── node
    │   └── node.go
    ├── scheduler
    │   └── scheduler.go
    ├── task
    │   └── task.go
    └── worker
        └── worker.go
```

* Call docker engine 
```
curl --unix-socket /var/run/docker.sock http://v1.45/containers/${container-id}/json
```