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

### 3. 使用说明
1. 运行后端

    go run test.go

2. 访问localhost:8080/index/
3. 点击出现的Send按钮，user.tpl会将ajax_data中的内容发送给后端，后端在验证后会将数据存储在数据库中
4. 若前端提示Success说明成功，Over Time说明时间戳与服务器时间相差过大（20s），Sign Dismatch说明Sign值不匹配，相应错误信息会在后端提示