# goUtil

一些常用工具类的开发包

- ## [CA 证书生成及解析](https://github.com/carmel/goUtil/ceti)
  证书相关的文件有如下格式：
  - **`.key`** - 通常指私钥
  - **`.csr`** - 证书签名请求 ( `Certificate Signing Request` )。可理解成公钥，生成证书时需将其发送给权威的证书颁发机构
  - **`.crt`** - 证书 ( `Certificate` )
  - **`.crl`** - 证书吊销列表 ( `Certificate Revocation List` )。 CRL 是由证书颁发机构 ( CA ) 维护的列表，并记录已撤销的证书的序列号。
  - **`.pem`** - `Privacy Enhanced Mail`，内容是以`BASE64`编码，并以`"-----BEGIN..."`开头, `"-----END..."`结尾的文本。`Apache`和`*NIX`系统倾向于使用这种编码格式。
  - **`.der`** - `Distinguished Encoding Rules`，不可读的二进制格式。

`X.509`是一种证书格式。对 X.509 证书来说，认证者总是 CA 或由 CA 指定的人，一份 X.509 证书是一些标准字段的集合，这些字段包含有关用户或设备及其相应公钥的信息。

`X.509`的证书文件一般以`.crt`结尾，根据该文件的内容编码格式，也可有以下两种格式：

**注**：该库生成的公钥格式均为`pem`，私钥格式均为`key`。

- ## [队列](https://github.com/carmel/goUtil/deque)

- ## [线程池](https://github.com/carmel/goUtil/pool)
