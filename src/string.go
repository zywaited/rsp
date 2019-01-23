package main

import (
    "reflect"
)

func (p *Rsp) Set(key string, value interface{}) bool {
    retval, _ := doWithKey(p, reflect.Bool, "SET", key, value)
    return retval.(bool)
}

func (p *Rsp) Get(key string) string {
    retval, _ := doWithKey(p, reflect.String, "GET", key)
    return retval.(string)
}

func (p *Rsp) Exists(key string) bool {
    retval, _ := doWithKey(p, reflect.Bool, "EXISTS", key)
    return retval.(bool)
}

func (p *Rsp) Expire(key string, timeout int) bool {
    retval, _ := doWithKey(p, reflect.Bool, "EXPIRE", key, timeout)
    return retval.(bool)
}

func (p *Rsp) ExpireAt(key string, timestamp int64) bool {
    retval, _ := doWithKey(p, reflect.Bool, "EXPIREAT", key, timestamp)
    return retval.(bool)
}

func (p *Rsp) Incr(key string, num int64) bool {
    retval, _ := doWithKey(p, reflect.Bool, "INCR", key, num)
    return retval.(bool)
}

func (p *Rsp) Decr(key string, num int64) bool {
    retval, _ := doWithKey(p, reflect.Bool, "DECR", key, num)
    return retval.(bool)
}

func (p *Rsp) Delete(key string) bool {
    retval, _ := doWithKey(p, reflect.Bool, "DEL", key)
    return retval.(bool)
}

