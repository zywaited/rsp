### rsp

PHP REDIS 扩展（使用GO开发）

### 现状

* 与 phpredis扩展调用相似
* 支持连接池
* 支持GO协程和管道

### 环境

* Linux/Unix system
* PHP 5.5+/7.x
* go version 1.9+

### 依赖
1 glide: 详情参考glide.yaml
2 <https://github.com/zywaited/php-go>: <https://github.com/kitech/php-go>
  2.1 增加链式操作（$this）
  2.2 增加interface{}返回
  2.3 增加nil返回
  2.4 修复float报错
3 <https://github.com/garyburd/redigo>

### 编译和安装

```sh
git clone https://github.com/zywaited/rsp

glide install

# pull request已提交
cd rsp/vendor/github.com
ln -s zywaited kitech

cd rsp/src
# 修改Makefile中PHPCFG为自己的php-config路径
make

# 启动redis（6379），无密码
php -d extension=rsp/src/rsp.so rsp/examples/test.php
```

### 实例
用法参考: rsp/examples/test.php
