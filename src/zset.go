package main

import (
    "reflect"
    "errors"
)

func (p *Rsp) ZRange(key string, start, end int64) []string {
    retval, _ := doWithKey(p, STRINGS, "ZRANGE", key, start, end)
    return retval.([]string)
}

func (p *Rsp) ZRangeWithScore(key string, start, end int64) map[string]float64 {
    retval, _ := doWithKey(p, FLOATMAP, "ZRANGE", key, start, end, "WITHSCORES")
    return retval.(map[string]float64)
}

func (p *Rsp) ZAdd(key string, score float64, value interface{}) bool {
    retval, _ := doWithKey(p, reflect.Bool, "ZADD", key, score, value)
    return retval.(bool)
}

func parseZAddArgs(items []interface{}) ([]interface{}, error) {
    var args []interface{}
    for i, v := range items {
        if (i + 1) % 2 == 0 {
            args = append(args, v)
            continue
        }

        switch reflect.TypeOf(v).Kind() {
            case reflect.Float64:
                fallthrough
            case reflect.Float32:
                fallthrough
            case reflect.Int64:
                fallthrough
            case reflect.Uint64:
                fallthrough
            case reflect.Int32:
                fallthrough
            case reflect.Uint32:
                fallthrough
            case reflect.Int:
                fallthrough
            case reflect.Uint:
                fallthrough
            case reflect.Int16:
                fallthrough
            case reflect.Uint16:
                fallthrough
            case reflect.Int8:
                fallthrough
            case reflect.Uint8:
                args = append(args, v)
            default:
                err := errors.New("rsp: zadd score must be numric")
                return nil, err
        }
    }

    return args, nil
}

func (p *Rsp) ZRem(key string, value interface{}) bool {
    retval, _ := doWithKey(p, reflect.Bool, "ZREM", key, value)
    return retval.(bool)
}

func (p *Rsp) ZCard(key string) int64 {
    retval, _ := doWithKey(p, reflect.Int64, "ZCARD", key)
    return retval.(int64)
}

