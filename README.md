# rdb-viewer

redis rdb文件解析工具 


## 安装
```
go get -u github.com/hidu/rdb-viewer
```

## 使用

### 基本用法
```
$ rdb-viewer dump.rdb
hset    "xxxxx1" 6
hset    "yyyyy2" 7
```

输出 3列: [string "abc" 12] 对应为 [类型 key 内容长度],eg  


* 只输出string类型的 *：  
```
 $ rdb_viewer -types string part1.rdb
```

输出string和set，而且输出具体内容  
```
$ rdb_viewer -types string,set -val part1.rdb
string  "xxx:phone:xxx:code"      7       value:  "verifed"       expiry: 0
string  "xxx:sms:counter"  1       value:  "1"     expiry: 0
```

如上 第三列以后是 key： value 格式  

### 输出json，使用其他脚本处理

```
rdb_viewer -json -val hd.rdb|php script/parse.php 
```


`script/parse.php` 的内容大致如下
```php
<?php
while(!feof(STDIN)){
    $line=fgets(STDIN);
    $obj=json_decode($line,true);
    if(is_array($obj)){
        foreach ($obj as $k=>$v){
            if(strpos($k,"_b")){
                $obj[substr($k, 0,strlen($k)-2)]=base64_decode($v);
                unset($obj[$k]);
            }
        }
    }
    //your code
    print_r($obj);
}
```
注：redis的数据是二进制的，key,value,member,field都是二进制的，输出的json内容中的字段均以`_b`结尾，如`key_b`
`$key`=base64_decode(`$key_b`)


