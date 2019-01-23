package main

import (
    "reflect"
)

func (p *Rsp) LPush(key string, value interface{}) int64 {
    retval, _ := doWithKey(p, reflect.Int64, "LPUSH", key, value)
    return retval.(int64)
}

func (p *Rsp) LPop(key string) string {
    retval, _ := doWithKey(p, reflect.String, "LPOP", key)
    return retval.(string)
}

func (p *Rsp) RPush(key string, value interface{}) int64 {
    retval, _ := doWithKey(p, reflect.Int64, "RPUSH", key, value)
    return retval.(int64)
}

func (p *Rsp) RPop(key string) string {
    retval, _ := doWithKey(p, reflect.String, "RPOP", key)
    return retval.(string)
}

func (p *Rsp) LRange(key string, start, end int64) []string {
    retval, _ := doWithKey(p, STRINGS, "LRANGE", key, start, end)
    return retval.([]string)
}

func (p *Rsp) LSet(key string, index int64, value interface{}) bool {
    retval, _ := doWithKey(p, reflect.Bool, "LSET", key, index, value)
    return retval.(bool)
}

func (p *Rsp) LLen(key string) int64 {
    retval, _ := doWithKey(p, reflect.Int64, "LLEN", key)
    return retval.(int64)
}

