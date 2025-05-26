# 资料站系统

一个基于Go Gin框架开发的资料管理系统，提供资料展示、下载和管理功能。

## 功能特点

### 公共功能（端口8080）
- 资料列表展示
- 资料下载

### 管理功能（端口8081）
- 管理员登录
- 资料上传
- 资料删除

## 技术栈

### 后端
- Go 1.21+
- Gin Web框架
- MySQL 数据库
- Redis 缓存（可选）
- JWT 认证

### 前端
- 原生HTML/CSS/JavaScript
- 响应式设计
- 简约现代UI

## 系统要求

- Go 1.21 或更高版本
- MySQL 5.7 或更高版本
- Redis 6.0 或更高版本（可选）
- 现代浏览器（Chrome、Firefox、Safari等）

## 快速开始

### 1. 克隆项目

```bash
git clone [项目地址]
cd servermanager
```

### 2. 配置数据库

1. 创建MySQL数据库：
```sql
CREATE DATABASE servermanager CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

2. 修改配置文件 `backend/config/config.yaml`：
```yaml
database:
  mysql:
    host: "localhost"
    port: 3306
    user: "your-username"
    password: "your-password"
    dbname: "servermanager"
```

### 3. 安装依赖

```bash
cd backend
go mod tidy
```

### 4. 启动服务

```bash
go run main.go
```

服务将在以下端口启动：
- 公共服务：http://localhost:8080
- 管理服务：http://localhost:8081

## 使用说明

### 访问公共页面
1. 打开浏览器访问 http://localhost:8080
2. 浏览资料列表
3. 点击资料卡片下载文件

### 访问管理页面
1. 打开浏览器访问 http://localhost:8081
2. 使用默认密码登录：admin123
3. 上传新资料：
   - 填写资料名称
   - 选择文件
   - 点击上传
4. 删除资料：
   - 在资料列表中找到要删除的资料
   - 点击删除按钮
   - 确认删除

## 项目结构

```
servermanager/
├── backend/
│   ├── config/
│   │   ├── config.yaml
│   │   ├── database.go
│   │   └── redis.go
│   ├── internal/
│   │   ├── handlers/
│   │   │   ├── admin.go
│   │   │   └── public.go
│   │   ├── middleware/
│   │   │   └── auth.go
│   │   └── models/
│   │       └── material.go
│   ├── uploads/
│   ├── go.mod
│   └── main.go
└── frontend/
    ├── index.html
    ├── admin.html
    ├── styles.css
    ├── script.js
    └── admin.js
```

## 安全说明

1. 生产环境部署前请修改以下配置：
   - 修改管理员密码
   - 修改JWT密钥
   - 配置HTTPS
   - 设置适当的文件上传限制

2. 建议的安全措施：
   - 使用环境变量存储敏感信息
   - 实现请求速率限制
   - 添加文件类型验证
   - 实现日志记录
   - 定期备份数据

## 常见问题

1. 数据库连接失败
   - 检查MySQL服务是否运行
   - 验证数据库配置是否正确
   - 确认数据库用户权限

2. 文件上传失败
   - 检查上传目录权限
   - 确认文件大小是否超限
   - 验证文件类型是否允许

3. 管理员登录失败
   - 确认密码是否正确
   - 检查JWT配置
   - 查看服务器日志

## 开发计划

- [ ] 添加资料分类功能
- [ ] 实现资料搜索
- [ ] 添加用户管理
- [ ] 支持文件预览
- [ ] 添加数据统计
- [ ] 优化文件存储

## 贡献指南

1. Fork 项目
2. 创建特性分支
3. 提交更改
4. 推送到分支
5. 创建 Pull Request

## 许可证

MIT License 