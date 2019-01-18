package main

import (
    "flag"
    "github.com/DGHeroin/etcd-go-discovery/discovery"
    "time"
)

var (
    name = flag.String("name", "pod-app", "node name")
)

func main()  {
    flag.Parse()
    data := []byte("hello")
    endpoints := []string{"http://127.0.0.1:2379"}

    pod, err := discovery.NewPod(*name, "pods", data,
        discovery.WithEndpoints(endpoints))
    if err != nil {
        panic(err)
    }
    time.AfterFunc(time.Second * 5, func() {
        pod.Stop()
    })
    pod.Serve()
}