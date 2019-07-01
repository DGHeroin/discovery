package discovery

import (
    "github.com/coreos/etcd/clientv3"
    "log"
    "context"
    "fmt"
    "encoding/base64"
    "sync/atomic"
)

type Master interface {
    Serve()
    Count() int64
    HandleFunc(cb func(EventType, string, []byte))
    Stop()
}

type master struct {
    prefix  string
    client  *clientv3.Client
    count   int64
    cb      func(EventType, string, []byte)
    stop    chan error
    options Options
}

type EventType int

func (e EventType) String() string {
    if e == 0 {
        return "ADD"
    } else if e == 1 {
        return "DELETE"
    }
    return "Unknown EventType"
}

const (
    EventTypeAdd    = EventType(0)
    EventTypeDelete = EventType(1)
)

func NewMaster(prefix string, opts ...OptionFunc) (Master, error) {
    options := defaultOptions
    for _, o := range opts {
        o(&options)
    }

    cli, err := clientv3.New(clientv3.Config{
        Endpoints:   options.Endpoints,
        Username:    options.Username,
        Password:    options.Password,
        DialTimeout: options.DialTimeout,
    })
    if err != nil {
        log.Println(err)
        return nil, err
    }

    master := &master{
        prefix: prefix,
        client: cli,
        stop:   make(chan error),
        options:options,
    }

    return master, nil
}

func (m *master) Serve() {
    key := fmt.Sprintf("/%v/", m.prefix)
    if m.options.debug {
        log.Printf("master monitoring [%v]", key)
    }

    rch := m.client.Watch(context.Background(), key, clientv3.WithPrefix())
StopMaster:
    for {
        select {
        case <-m.stop:
            log.Printf("Stop master")
            break StopMaster
        case resp := <-rch:
            m.onResponse(resp)
        }
    }
}

func (m *master) onResponse(resp clientv3.WatchResponse) {
    for _, ev := range resp.Events {
        switch ev.Type {
        case clientv3.EventTypePut:
            key := string(ev.Kv.Key)
            value := string(ev.Kv.Value)
            var data []byte
            if val, err := base64.StdEncoding.DecodeString(value); err != nil {
                log.Println(err)
                break
            } else {
                data = val
            }
            atomic.AddInt64(&m.count, 1)

            if m.cb != nil {
                m.cb(EventTypeAdd, key, data)
            }

        case clientv3.EventTypeDelete:
            key := string(ev.Kv.Key)
            atomic.AddInt64(&m.count, -1)

            if m.cb != nil {
                m.cb(EventTypeDelete, key, nil)
            }
        }
    }

}

func (m *master) Count() int64 {
    return m.count
}
func (m *master) HandleFunc(cb func(EventType, string, []byte)) {
    m.cb = cb
}

func (m *master) Stop() {
    m.stop <- nil
}
