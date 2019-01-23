package main

import (
    "github.com/garyburd/redigo/redis"
    "strings"
    "log"
    "os"
    "reflect"
    "errors"
    "strconv"
)

const (
    INVALID_CMD int = iota
    MULTI_CMD
    PIPE_CMD
)

// redis特定类型
const (
    FLOATMAP = reflect.UnsafePointer + 10
    STRINGMAP = FLOATMAP + 1
    STRINGS = STRINGMAP + 1
)

var std = log.New(os.Stdout, "", log.LstdFlags)

func float64Map(result interface{}, err error) (map[string]float64, error) {
    values, err := redis.Values(result, err)
    if err != nil {
        return nil, err
    }

    vLen := len(values)
    if vLen % 2 != 0 {
        return nil, errors.New("rsp: float64Map expects even number of values result")
    }

    m := make(map[string]float64, vLen / 2)
    for i := 0; i < vLen; i += 2 {
        key, okKey := values[i].([]byte)
        value, okValue := values[i+1].([]byte)
        if !okKey || !okValue {
            return nil, errors.New("rsp: flocat64Map key not a bulk float64 value")
        }

        f, err := strconv.ParseFloat(string(value), 64)
        if err != nil {
            return nil, err
        }

        m[string(key)] = f
    }

    return m, nil
}

func originDo(p *Rsp, t reflect.Kind, cmd string, args ...interface{}) (retval interface{}, err error) {
    defer func() {
        if err := recover(); err != nil {
            std.Println(err)
        }
    }()

    if p.rp == nil {
        p.createRp()
    }

    var r redis.Conn
    if p.rx != nil {
        r = *(p.rx)
    } else {
        r = p.rp.Get()
    }

    cmd = strings.ToUpper(cmd)
    switch p.multi {
        case PIPE_CMD:
            if p.rx == nil {
                p.rx = &r
            }

            if cmd == "" {
                defer r.Close()
                p.rx = nil
                _, err = r.Do("", args...)
            } else {
                err = r.Send(cmd, args...)
            }

            retval, err = parseMulti(t, err, err == nil)
        case MULTI_CMD:
            fallthrough
        default:
            var pipe bool = false
            if cmd == "EXEC" || cmd == "MULTI" {
                if cmd == "EXEC" {
                    defer r.Close()
                    p.rx = nil
                } else if p.rx == nil {
                    p.rx = &r
                }

                pipe = true
            }

            if !pipe && p.rx != nil {
                pipe = true
            }

            retval, err = r.Do(cmd, args...)
            if pipe {
                retval, err = parseMulti(t, err, err == nil)
            } else {
                switch t {
                    case reflect.Bool:
                        retval, err = redis.Bool(retval, err)
                    case reflect.String:
                        retval, err = redis.String(retval, err)
                    case reflect.Int:
                        retval, err = redis.Int(retval, err)
                    case STRINGMAP:
                        retval, err = redis.StringMap(retval, err)
                    case reflect.Int64:
                        retval, err = redis.Int64(retval, err)
                    case STRINGS:
                        retval, err = redis.Strings(retval, err)
                    case FLOATMAP:
                        retval, err = float64Map(retval, err)
                    default:
                        err = errors.New("the type is not known")
                }
            }
    }

    saveLog(err)
    return
}

func do(p *Rsp, t reflect.Kind, cmd string, args ...interface{}) (interface{}, error) {
    if p.isSync {
        p.mu.Lock()
        if p.syncNum >= p.maxActive {
            err := errors.New("arrive to max chan num")
            saveLog(err)
            p.mu.Unlock()
            return parseMulti(t, err, false)
        }

        p.syncNum++
        p.mu.Unlock()
        go func () {
            retval, _ := originDo(p, t, cmd, args...)
            p.data <- retval
            p.mu.Lock()
            p.readyNum++
            p.mu.Unlock()
        }()

        return parseMulti(t, nil, true)
    }

    return originDo(p, t, cmd, args...)
}

func doWithKey(p *Rsp, t reflect.Kind, cmd string, key string, args ...interface{}) (interface{}, error) {
    // newArgs := make([]interface{}, len(args) + 1)
    var newArgs []interface{}
    if p.prefix != "" {
        // 数据量少直接拼接,不用buffer
        key = p.prefix + key
    }

    newArgs = append(newArgs, key)
    newArgs = append(newArgs, args...)
    return do(p, t, cmd, newArgs...)
}

func saveLog(err error) {
    if err != nil {
        std.Println(err)
        // 不中断
        // panic(err.Error())
    }
}

func auth(r redis.Conn, auth string) bool {
    retval, err := redis.Bool(r.Do("AUTH", auth))
    saveLog(err)
    return retval
}

func selectDb(r redis.Conn, db int) bool {
    retval, err := redis.Bool(r.Do("SELECT", db))
    saveLog(err)
    return retval
}

func (p *Rsp) Auth(auth string) bool {
    retval, _ := originDo(p, reflect.Bool, "AUTH", auth)
    return retval.(bool)
}

func (p *Rsp) Select(db int) bool {
    retval, _ := originDo(p, reflect.Bool, "SELECT", db)
    return retval.(bool)
}


func (p *Rsp) Multi(cmd int) bool {
    switch cmd {
        case MULTI_CMD:
            p.multi = MULTI_CMD
            retval, _ := do(p, reflect.Bool, "MULTI")
            return retval.(bool)
        case PIPE_CMD:
            fallthrough
        default:
            p.multi = PIPE_CMD
    }

    return true
}

func (p *Rsp) Exec() bool {
    var cmd string = ""
    if p.multi == MULTI_CMD {
        cmd = "EXEC"
    }

    retval, _ := do(p, reflect.Bool, cmd)
    return retval.(bool)
}

func parseMulti(t reflect.Kind, e error, status bool) (retval interface{}, err error) {
    err = e
    switch t {
        case reflect.Bool:
            retval = status
        case reflect.String:
            if status {
                retval = "1"
            } else {
                retval = "0"
            }
        case reflect.Int:
            var rtv int
            if status {
                rtv = 1
            } else {
                rtv = 0
            }

            retval = rtv
        case reflect.Int64:
            var rtv int64
            if status {
                rtv = 1
            } else {
                rtv = 0
            }

            retval = rtv
        case STRINGMAP:
            rtv := make(map[string]string, 1)
            if status {
                 rtv["status"] = "1"
            } else {
                 rtv["status"] = "0"
            }

             retval = rtv
        case STRINGS:
            rtv := make([]string, 1)
            if status {
                 rtv = append(rtv, "1")
            } else {
                 rtv = append(rtv, "0")
            }

            retval = rtv
        case FLOATMAP:
            rtv := make(map[string]float64, 1)
            if status {
                rtv["status"] = 1
            } else {
                rtv["status"] = 0
            }

            retval = rtv
        default:
            err = errors.New("the type is not known")
    }

    return
}

func (p *Rsp) GetData() interface{} {
    p.mu.Lock()
    if p.syncNum < 1 {
        p.mu.Unlock()
        return nil
    }

    p.syncNum--
    p.readyNum--
    p.mu.Unlock()
    return <- p.data
}

func (p *Rsp) IsAllDataReady() bool {
    return p.syncNum == p.readyNum
}

func (p *Rsp) hasDataReady() bool {
    return p.readyNum > 0
}

