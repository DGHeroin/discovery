package discovery

import "time"

type Options struct {
    Endpoints   []string      `json:"endpoints"`
    Username    string        `json:"username"`
    Password    string        `json:"password"`
    DialTimeout time.Duration `json:"dial-timeout"`
    TTL         int64         `json:"ttl"`
    debug       bool
}

var defaultOptions = Options{
    Endpoints: []string{"http://127.0.0.1:2379"},
    Username:  "",
    Password:  "",
    TTL:       5,
    debug:     false,
}

type OptionFunc func(options *Options)

func WithEndpoints(endpoints []string) OptionFunc {
    return func(o *Options) {
        o.Endpoints = endpoints
    }
}

func WithUsername(username string) OptionFunc {
    return func(o *Options) {
        o.Username = username
    }
}

func WithPassword(password string) OptionFunc {
    return func(o *Options) {
        o.Password = password
    }
}

func WithDialTimeout(timeout time.Duration) OptionFunc {
    return func(o *Options) {
        o.DialTimeout = timeout
    }
}

func WithDebug(debug bool) OptionFunc {
    return func(o *Options) {
        o.debug = debug
    }
}

func WithTTL(ttl int64) OptionFunc {
    return func(o *Options) {
        o.TTL = ttl
    }
}
