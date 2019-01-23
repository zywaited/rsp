package main

import (
    "github.com/garyburd/redigo/redis"
    "errors"
    "strconv"
    "time"
    "sync"
)

type Rsp struct {
    tcp bool
    host string
    port int
    unix string
    auth string
    dbIndex int
    prefix string
    multi int
    rp *redis.Pool
    data chan interface{}
    rx *redis.Conn
    maxIdle int
    maxActive int
    timeout time.Duration
    maxThreads int
    isSync bool
    syncNum int
    readyNum int
    mu sync.Mutex
}

func NewRsp() *Rsp {
    return &Rsp{
        tcp: false,
        auth: "",
        dbIndex: -1,
        prefix: "",
        multi: MULTI_CMD,
        rp: nil,
        maxIdle: 30,
        maxActive: 10,
        maxThreads: 10,
        isSync: false,
    }
}

func (p *Rsp) SetTcp(host string, port int) *Rsp {
    p.tcp = true
    p.host = host
    p.port = port
    return p;
}

func (p *Rsp) SetUnix(unix string) *Rsp {
    p.tcp = false
    p.unix = unix
    return p
}

func (p *Rsp) SetRedisConfig(auth string, db int) *Rsp {
    p.auth = auth
    p.dbIndex = db
    return p
}

func (p *Rsp) SetRedisPrefixKey(prefix string) *Rsp {
    p.prefix = prefix
    return p
}

func (p *Rsp) initRedis(r redis.Conn) (redis.Conn, error) {
    if p.auth != "" {
        if !auth(r, p.auth) {
            r.Close();
            return nil, errors.New("auth redis error")
        }
    }

    if p.dbIndex > 0 {
        if !selectDb(r, p.dbIndex) {
            r.Close();
            return nil, errors.New("select redis's db error")
        }
    }

    return r, nil
}

func (p *Rsp) getUnixDial() (redis.Conn, error) {
    r, err := redis.Dial("unix", p.unix)
    if err != nil {
        return nil, err
    }

    return p.initRedis(r)
}

func (p *Rsp) getTcpDial() (redis.Conn, error) {
    r, err := redis.Dial("tcp", p.host + ":" + strconv.Itoa(p.port))
    if err != nil {
        return nil, err
    }

    return p.initRedis(r)
}

func (p *Rsp) createRp() (*Rsp, error) {
    if p.tcp {
        if p.host == "" || p.port <= 0 {
            return nil, errors.New("redis's host or port not valid")
        }

        // default maxidle 30
        p.rp = &redis.Pool{
            Dial: func () (redis.Conn, error) {
                return p.getTcpDial()
            },
            MaxIdle: p.maxIdle,
            MaxActive: p.maxActive,
            IdleTimeout: p.timeout,
            Wait: true,
        }

        return p, nil
    }

    if p.unix == "" {
        return nil, errors.New("redis's sock not valid")
    }

    p.rp = &redis.Pool{
        Dial: func () (redis.Conn, error) {
          return p.getUnixDial()
        },  
        MaxIdle: p.maxIdle,
        MaxActive: p.maxActive,
        IdleTimeout: p.timeout,
        Wait: true,
    }

    return p, nil
}

func (p *Rsp) SetUp() *Rsp {
    p.createRp()
    p.data = make(chan interface{}, p.maxActive)
    p.syncNum = 0
    p.readyNum = 0
    return p
}

// set pool config
func (p *Rsp) SetPoolConfig(maxIdle, maxActive int, timeout time.Duration) *Rsp {
    p.maxIdle = maxIdle
    p.maxActive = maxActive
    p.timeout = timeout * time.Second
    return p
}

func (p *Rsp) SetMaxThreads(maxThreads int) *Rsp {
    if maxThreads <= 0 {
        return p
    }

    p.maxThreads = maxThreads
    return p
}

func (p *Rsp) SetSync(isSync bool) *Rsp {
    p.isSync = isSync
    return p
}

// get pool status
func (p *Rsp) GetActiveCount() int {
    if p.rp == nil {
        p.createRp()
    }

    return p.rp.ActiveCount()
}

func (p *Rsp) GetIdleCount() int {
    if p.rp == nil {
        p.createRp()
    }

    return p.rp.IdleCount()
}

