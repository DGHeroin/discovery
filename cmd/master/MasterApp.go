package main

import (
    "github.com/DGHeroin/etcd-go-discovery/discovery"
    "log"
    "time"
)

func main()  {
    endpoints := []string{"http://127.0.0.1:2379"}
    var master discovery.Master
    if m, err := discovery.NewMaster( "pods", discovery.WithEndpoints(endpoints)); err != nil {
        panic(err)
    } else {
        master = m
    }

    master.HandleFunc(func(eventType discovery.EventType, key string, value []byte) {
        log.Printf("Event:%-6s|%-20v|%-20v|Pod Num=%v",
                eventType, key, string(value),
                master.Count())
    })

    time.AfterFunc(time.Minute*5, func() {
        master.Stop()
    })

    master.Serve()
}
