## go 微服务网关项目后端

### 开源地址

码云：`https://gitee.com/bclz_xc/go-gateway-service`

github：`https://github.com/lilj01/go-gateway-service`

### 涉及技术说明

- redis 缓存
- mysql 数据库
- golang 语言
- gin 、gorm 框架
- gin_scaffold 脚手架 （地址:`https://github.com/e421083458/gin_scaffold`）

### 本地运行

```
运行main.go
1.ide运行main函数
2.或命令行 go run main.go
```

### 本地运行 swagger 地址

`http://127.0.0.1:8880/swagger/index.html`

### 服务器部署

待补充，会补充 k8s 部署方式

## go 微服务网关项目前端地址

### 开源地址

码云:`https://gitee.com/bclz_xc/go-gateway-view`

github: 待迁移

### 涉及技术说明

vue-element-admin，文档地址：

`https://panjiachen.gitee.io/vue-element-admin-site/zh/`

本地运行

```apl
# 安装依赖
npm install
或
npm install --registry=https://registry.npm.taobao.org

# 本地开发 启动项目
npm run serve
```

## 目前已完成的功能

### 管理员

管理员登录，退出，基本信息，密码修改



### 服务管理

服务列表，http服务添加（前台未完成）







## 功能展示

### 登录

![](http://bclz_xc.gitee.io/lilj_01-static/go/go_gateway/login_page.png)





### 服务列表

![](https://bclz_xc.gitee.io/lilj_01-static/go/go_gateway/service_list_page.png)