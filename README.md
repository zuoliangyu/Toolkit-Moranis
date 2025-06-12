# 资料站系统

一个基于Go Gin框架开发的资料管理系统，提供资料展示、下载和管理功能。

## 功能特点

### 公共功能（端口8080）
- 📁 文件夹树形展示，支持多层级目录结构（最多99层）
- 📄 资料列表展示，支持按文件夹筛选
- 📥 资料下载
- 🔍 文件夹导航和资料查找

### 管理功能（端口8081）
- 🔐 管理员登录认证
- 📁 文件夹管理：创建、重命名、删除文件夹
- 📤 资料上传，支持选择目标文件夹
- 🗑️ 资料删除
- 🌳 文件夹树形管理界面

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
2. 使用文件夹功能：
   - 在左侧文件夹树中浏览文件夹结构
   - 点击文件夹名称查看该文件夹下的资料
   - 点击展开/折叠按钮（▶/▼）管理文件夹显示
   - 点击"所有资料"查看全部资料
3. 下载资料：
   - 在资料列表中找到需要的文件
   - 点击下载链接即可下载

### 访问管理页面
1. 打开浏览器访问 http://localhost:8081
2. 使用默认密码登录：admin123
3. 文件夹管理：
   - **创建文件夹**：
     - 在"文件夹管理"区域输入文件夹名称
     - 选择父文件夹（可选，默认为根目录）
     - 点击"创建文件夹"按钮
   - **重命名文件夹**：
     - 在文件夹列表中找到要重命名的文件夹
     - 点击"重命名"按钮
     - 输入新名称并确认
   - **删除文件夹**：
     - 确保文件夹为空（无子文件夹和资料文件）
     - 点击"删除"按钮并确认
4. 上传新资料：
   - 填写资料名称
   - 选择目标文件夹（可选，默认为根目录）
   - 选择文件
   - 点击上传
5. 删除资料：
   - 在资料列表中找到要删除的资料
   - 点击删除按钮
   - 确认删除

## 项目结构

```
servermanager/
├── backend/
│   ├── config/
│   │   ├── config.yaml          # 配置文件
│   │   ├── database.go          # 数据库配置
│   │   └── redis.go             # Redis配置
│   ├── internal/
│   │   ├── handlers/
│   │   │   ├── admin.go         # 管理员API处理器
│   │   │   └── public.go        # 公共API处理器
│   │   ├── middleware/
│   │   │   └── auth.go          # JWT认证中间件
│   │   └── models/
│   │       ├── folder.go        # 文件夹数据模型
│   │       └── material.go      # 资料数据模型
│   ├── uploads/                 # 文件上传目录
│   ├── go.mod
│   └── main.go                  # 程序入口
└── frontend/
    ├── index.html               # 公共页面
    ├── admin.html               # 管理员页面
    ├── styles.css               # 样式文件
    ├── script.js                # 公共页面脚本
    └── admin.js                 # 管理员页面脚本
```

## API 接口文档

### 公共API（端口8080）

#### 资料相关
- `GET /api/materials` - 获取资料列表，支持`folder_id`参数筛选
- `GET /api/materials/:id/download` - 下载指定资料

#### 文件夹相关
- `GET /api/folders` - 获取文件夹树形结构
- `GET /api/folders/:id/materials` - 获取指定文件夹下的资料

### 管理员API（端口8081，需要JWT认证）

#### 认证相关
- `POST /api/login` - 管理员登录
- `POST /api/admin/login` - 管理员登录（备用）

#### 资料管理
- `POST /api/materials` - 上传资料，支持`folder_id`参数
- `DELETE /api/materials/:id` - 删除指定资料

#### 文件夹管理
- `GET /api/folders` - 获取文件夹列表（管理员视图）
- `POST /api/folders` - 创建文件夹
- `PUT /api/folders/:id` - 重命名文件夹
- `DELETE /api/folders/:id` - 删除文件夹（需为空文件夹）

### 请求示例

#### 创建文件夹
```bash
curl -X POST http://localhost:8081/api/folders \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "新文件夹",
    "parent_id": 1,
    "description": "文件夹描述"
  }'
```

#### 上传资料到指定文件夹
```bash
curl -X POST http://localhost:8081/api/materials \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -F "name=资料名称" \
  -F "folder_id=1" \
  -F "file=@/path/to/file.pdf"
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

4. 文件夹操作问题
   - **无法创建文件夹**：检查文件夹名称是否重复，确认层级未超过99层
   - **无法删除文件夹**：确保文件夹为空（无子文件夹和资料文件）
   - **文件夹树不显示**：检查网络连接，查看浏览器控制台错误信息
   - **文件夹层级混乱**：刷新页面或重新登录管理员界面

5. 文件上传到文件夹失败
   - 确认目标文件夹存在且有效
   - 检查文件夹ID是否正确
   - 验证管理员权限

## 开发计划

### 已完成功能 ✅
- [x] 文件夹管理功能（多层级目录结构，最多99层）
- [x] 文件夹树形展示和导航
- [x] 资料按文件夹分类和筛选
- [x] 管理员文件夹CRUD操作
- [x] 响应式设计和移动端适配

### 计划中功能 📋
- [ ] 实现资料搜索功能
- [ ] 添加用户管理系统
- [ ] 支持文件预览（PDF、图片等）
- [ ] 添加数据统计和分析
- [ ] 优化文件存储（云存储支持）
- [ ] 实现文件夹权限管理
- [ ] 添加操作日志记录
- [ ] 支持批量文件操作
- [ ] 实现文件版本管理
- [ ] 添加文件标签系统

## 贡献指南

1. Fork 项目
2. 创建特性分支
3. 提交更改
4. 推送到分支
5. 创建 Pull Request

## 许可证

MIT License 