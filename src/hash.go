package main

import (
    "reflect"
)

// hash
func (p *Rsp) HSet(key, hashKey string, value interface{}) int {
    retval, _ := doWithKey(p, reflect.Int, "HSET", key, hashKey, value)
    return retval.(int)
}

func (p *Rsp) HGet(key, hashKey string) string {
    retval, _ := doWithKey(p, reflect.String, "HGET", key, hashKey)
    return retval.(string)
}

// hgetall 转成 map
func (p *Rsp) HGetAll(key string) map[string]string {
    retval, _ := doWithKey(p, STRINGMAP, "HGETALL", key)
    return retval.(map[string]string)
}

// 与php redis扩展保持一致 hmset使用数组
func (p *Rsp) HMSet(key string, hashKeys map[interface{}]interface{}) bool {
    // map to array(slice)
    if len(hashKeys) < 1 {
        return true
    }

    var args []interface{}
    for k, v := range hashKeys {
        args = append(args, k)
        args = append(args, v)
    }

    retval, _ := doWithKey(p, reflect.Bool, "HMSET", key, args...)
    return retval.(bool)
}

func (p *Rsp) HMGet(key string, hashKeys []string) map[string]string {
    retval, _ := doWithKey(p, STRINGMAP, "HGETALL", key)
    return retval.(map[string]string)
}

func (p *Rsp) HDel(key, hashKey string) bool {
    retval, _ := doWithKey(p, reflect.Bool, "HDEL", key, hashKey)
    return retval.(bool)
}

func (p *Rsp) HLen(key string) int64 {
    retval, _ := doWithKey(p, reflect.Int64, "HLEN", key)
    return retval.(int64)
}

func (p *Rsp) HExists(key, hashKey string) bool {
    retval, _ := doWithKey(p, reflect.Bool, "HEXISTS", key, hashKey)
    return retval.(bool)
}

func (p *Rsp) HIncrby(key, hashKey string, num int64) bool {
    retval, _ := doWithKey(p, reflect.Bool, "HINCRBY", key, hashKey, num)
    return retval.(bool)
}

