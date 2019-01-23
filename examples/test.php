<?php

use Waited\Rsp;

$p = new Rsp();
$p->setTcp('127.0.0.1', 6379); 
// 设置auth和db
// $p->setRedisConfig('', 1);

// 设置前缀
// $p->setRedisPrefixKey("Waited_");

// 设置启动
$p->setUp()->SetPoolConfig(5, 10, 2);

// // set get expire
// var_dump($p->set('waited', 'ok'));
// var_dump($p->get('waited'));
// var_dump($p->delete('waited', 30));
// 
// // hash
// $p->hSet('waited', 'test', 1);
// $p->hSet('waited', 'test_one', 2);
// var_dump($p->hGet('waited', 'test'));
// $p->hDel('waited', 'test');
// var_dump($p->hIncrby('waited', 'test_one', 1));
// var_dump($p->hGetAll('waited'));
// var_dump($p->hLen('waited'));
// var_dump($p->hExists('waited', 'test_one'));
// var_dump($p->Exists('waited'));
// var_dump($p->delete('waited', 30));
// 
// // list
// $p->lPush('waited', 1);
// $p->lPush('waited', 2);
// $p->rPush('waited', 3);
// $p->rPush('waited', 4);
// var_dump($p->lLen('waited'));
// var_dump($p->rPop('waited'));
// var_dump($p->lPop('waited'));
// $p->expireAt('waited', 1547642388);
// 
// // zset
// $p->zAdd('waited', 0, 1);
// $p->zAdd('waited', 1, 2);
// $p->zAdd('waited', 2, 3);
// $p->zAdd('waited', 3, 4);
// var_dump($p->zRem('waited', 2));
// var_dump($p->zCard('waited'));
// var_dump($p->zRange('waited', 0, 5));
// var_dump($p->zRangeWithScore('waited', 0, 5));
// var_dump($p->delete('waited', 30));

// pipe
// $p->multi(2);
// var_dump($p->multi(1));
// var_dump($p->lPush('waited', 1));
// var_dump($p->lPush('waited', 2));
// var_dump($p->rPop('waited'));
// var_dump($p->exec());

// chan sync
var_dump($p->getData());
$p->setSync(true);
// 写入顺序无法保证
$p->lPush('waited', 1);
$p->lPush('waited', 2);
while ($p->isAllDataReady());
var_dump($p->getData());
var_dump($p->getData());

$p->expire('waited', 30);
