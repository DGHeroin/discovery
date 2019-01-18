package discovery

import (
    "github.com/coreos/etcd/clientv3"
    "fmt"
    "context"
    "log"
    "encoding/base64"
)

type Pod interface {
    Serve() error
    Stop()
}

type pod struct {
    Name string
    stop    chan error
    leaseId clientv3.LeaseID
    client  *clientv3.Client
    prefix  string
    data    []byte
    options Options
}

func NewPod(name string, prefix string, data []byte, opts ...OptionFunc) (Pod, error) {
    options := defaultOptions
    for _, o := range opts {
        o(&options)
    }

    cli, err := clientv3.New(clientv3.Config{
        Endpoints: options.Endpoints,
        Username:options.Username,
        Password:options.Password,
        DialTimeout: options.DialTimeout,
    })
    if err != nil {
        return nil, err
    }
    return &pod{Name: name,
        data: data,
        stop: make(chan error),
        client: cli,
        prefix: prefix,
        options:options,
    },
    nil
}

func (p *pod) Serve() error {
    ch, err := p.keepAlive()
    if err != nil {
        log.Println(err)
        return err
    }
    for {
        select {
        case err := <-p.stop:
            p.revoke()
            return err
        case <-p.client.Ctx().Done():
            return fmt.Errorf("%v", "Server closed")
        case ka, ok := <-ch:
            if !ok {
                log.Println("keep alive channel closed")
                p.revoke()
                return nil
            }
            if p.options.debug {
                log.Printf("Recv reply form pod: %s, ttl:%d", p.Name, ka.TTL)
            }
        }
    }
}

func (p *pod) Stop() {
    p.stop <- nil
}

func (p *pod) keepAlive() (<-chan *clientv3.LeaseKeepAliveResponse, error) {
    key := fmt.Sprintf("/%s/%s", p.prefix, p.Name)
    if p.options.debug {
        log.Printf("pod write to [%v]", key)
    }
    value := base64.StdEncoding.EncodeToString(p.data)
    resp, err := p.client.Grant(context.TODO(), 10)
    if err != nil {
        return nil, err
    }
    _, err = p.client.Put(context.TODO(), key, value, clientv3.WithLease(resp.ID))
    if err != nil {
        log.Println(err)
        return nil, err
    }
    p.leaseId = resp.ID
    return p.client.KeepAlive(context.TODO(), resp.ID)
}

func (p *pod) revoke() error {
    _, err := p.client.Revoke(context.TODO(), p.leaseId)
    if err != nil {
        log.Println(err)
        return err
    }
    log.Printf("pod: %s stop", p.Name)
    return nil
}
