### GRPC
下载
```shell
# protoc compiler
https://github.com/protocolbuffers/protobuf/releases

# protoc go plugins
# protoc-gen-go
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
# protoc-gen-go-grpc
go install google.golang.org/protobuf/cmd/protoc-gen-go-grpc@latest
# protoc-gen-grpc-gateway
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
# protoc-gen-openapiv2
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest

# protoc typescript plugins

```

### WHEREIP
#### 纯真IP库
**纯真ip库获取及解密方式**：https://blog.dnomd343.top/qqwry.dat-analyse/#%E5%89%8D%E8%A8%80

**下载地址可能会变化，需要支持热更新**  
纯真ip库文件下载地址：https://gh-release.zu1k.com/HMBSbige/qqwry/qqwry.dat

**纯真ip库文件结构**  
分为三部分：
- 文件头：固定8字节，前4字节为索引区第一条索引偏移量，后4字节为索引区最后一条索引偏移量
- 记录区：索引后4字节指向的位置为终止ip，也就是**索引前4字节，与索引偏移量指向地址的后4字节组成一个ip区间**，随后的1个字节，如果为0x01或0x02，则记录数据需要再次读取3字节的偏移量进行跳转。
- 索引区：区域大小不定长，每条索引固定为7字节，前4字节记录起始ip，后3字节为区间在记录区的偏移量，索引区ip从小到大排列

索引区只记录了起始ip，从小到大排列，可以使用二分查找加速搜索速度，查找最后一个小于等于该ip的索引，对应`Leetcode:`[436. 寻找右区间](https://leetcode.cn/problems/find-right-interval/description/)。

记录模式：
- 0x01：后面3字节表示记录A记录B都需要偏移
- 0x02：后面3字节表示记录A需要偏移，记录B数据无需偏移，记录B读取到二进制0结束
- 其他：后面跟着记录数据，读取到二进制0结束

**0x01和0x02再次偏移后，仍有可能仍是0x01和0x02，需要递归处理。**

0x02模式下，3字节的偏移量一定代表了国家的名字

#### MAXMIND GEOIP库
下载地址：https://www.maxmind.com/en/accounts/1037613/geoip/downloads

**自带官方SDK库，免费查询支持到城市级别，但是解析内容全是英文，不适合国内用户。**