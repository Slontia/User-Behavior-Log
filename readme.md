# 接口说明

### 1. 数据格式
接收json格式数据，键值对规定如下：

| 键名 | 类型 | 说明 |
|:-:|:-:|:-:|
| time | int | 时间戳，（如1432710115）
| sign | string | 签名值（后面会介绍生成算法）
| action | string | 用户行为
| id | string | 用户id

### 2. 生成sign算法
1. 以 "id" + id + "action" + action + "time" + time 的形式进行字符串拼接
*如：id001actionlogintime1432710115*
2. 将上述字符串进行MD5加密，得到的结果就是sign
