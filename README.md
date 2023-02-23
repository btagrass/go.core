# Go核心库
## 简介
|包|名称|模块|
|:-:|:-:|:-:|
|app|应用|配置，日志|
|cmd|命令行|安装，卸载，启动，结束，状态|
|cmw|通用中间件|认证（Gin+Jwt+Casbin）|
|dao|泛型数据访问对象|获取，获取集合（列表），移除，保存，更新，事务|
|htp|超文本传输协议|Get，Post，PostFile，PostForm，SaveFile|
|mdl|泛型模型|模型接口，模型类|
|ntp|网络传输协议|UDP|
|r|请求结果|查询（query转map），输出Json|
|svc|泛型服务|数据库（DDL），缓存|
|sys|系统|部门管理，用户管理，角色管理，资源管理，字典管理，更新服务|
|trc|分时量比控制|广告、任务等量比控制|
|utl|实用方法|密码，执行，文件路径，编码，多媒体，网络，字符串集合|
## 安装
```bash
go get github.com/btagrass/go.core
```
