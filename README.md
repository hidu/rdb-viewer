# rdb-viewer
parse redis rdb file


## install
```
export GO15VENDOREXPERIMENT=1
go get -u github.com/hidu/rdb-viewer
```

## useage
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


