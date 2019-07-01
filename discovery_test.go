package discovery

import (
    "testing"
    "time"
)

var (
    endpoints = []string{"http://127.0.0.1:2379"}
)

func TestDiscovery(t *testing.T) {
    go startMaster(t)
    startPod(t)
}

func startMaster(t *testing.T)  {
    var err error
    var master Master
    master, err = NewMaster( "pods", WithEndpoints(endpoints));
    if err != nil {
        t.Fatal(err)
    }

    master.HandleFunc(func(eventType EventType, key string, value []byte) {
        t.Logf("Event:%-6s|%-20v|%-20v|Pod Num=%v",
            eventType, key, string(value),
            master.Count())
    })

    time.AfterFunc(time.Second * 5, func() {
        master.Stop()
    })

    master.Serve()
}

func startPod(t *testing.T) {
    data := []byte("hello-world")
    pod, err := NewPod("pod-name", "pods", data,
        WithEndpoints(endpoints))
    if err != nil {
        t.Fatal(err)
    }
    time.AfterFunc(time.Second * 2, func() {
        pod.Stop()
    })
    pod.Serve()
}
