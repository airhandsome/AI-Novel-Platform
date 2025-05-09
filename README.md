# AI Novel Platform

一个基于 Vue 3.0 和 Go Gin 的AI辅助小说创作平台，集成 DeepSeek AI 实现智能写作。

## 项目架构

### 前端技术栈
- Vue 3.0 + Vite
- Pinia (状态管理)
- Vue Router
- Axios (HTTP 请求)
- Element Plus (UI组件库)
- TailwindCSS (样式框架)

### 后端技术栈
- Gin Web Framework
- GORM (ORM框架)
- JWT (认证)
- Redis (缓存)
- Swagger (API文档)

### 数据库
- MySQL 8.0 (主数据库)
- Redis (缓存)

## 核心功能模块

### 1. 用户系统
- 用户注册/登录
- 个人信息管理
- 作品管理
- 收藏/订阅
- 阅读历史

### 2. 小说阅读系统
- 小说分类浏览
- 搜索功能
- 阅读器（支持多主题、字体大小调节）
- 书签/笔记
- 评论系统
- 推荐系统

### 3. AI写作系统
- DeepSeek集成
  - 写作辅助
  - 情节生成
  - 角色设计
  - 场景描写
  - 文风调整

### 4. 提示词系统
- 世界观构建
- 人物塑造
- 情节发展
- 场景描写
- 对话生成
- 文风调教

## 数据库设计

### 主要数据表
```sql
-- 用户表
CREATE TABLE users (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    username VARCHAR(50) UNIQUE,
    password_hash VARCHAR(255),
    email VARCHAR(100) UNIQUE,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);

-- 小说表
CREATE TABLE novels (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    title VARCHAR(200),
    author_id BIGINT,
    description TEXT,
    cover_url VARCHAR(255),
    category VARCHAR(50),
    status TINYINT,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);

-- 章节表
CREATE TABLE chapters (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    novel_id BIGINT,
    title VARCHAR(200),
    content TEXT,
    chapter_order INT,
    word_count INT,
    created_at TIMESTAMP
);

-- 提示词表
CREATE TABLE prompts (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    category VARCHAR(50),
    content TEXT,
    usage_count INT,
    created_at TIMESTAMP
);
```

## API接口规范 
/api/v1/
├── auth/
│ ├── login
│ ├── register
│ └── refresh-token
├── users/
│ ├── profile
│ ├── novels
│ └── favorites
├── novels/
│ ├── list
│ ├── detail
│ ├── chapters
│ └── search
└── ai/
├── generate
├── prompts
└── assist


## 开发环境搭建

### Docker 配置
```yaml
services:
  mysql:
    image: mysql:8.0
    environment:
      MYSQL_ROOT_PASSWORD: your_password
      MYSQL_DATABASE: novel_platform
    
  redis:
    image: redis:alpine
    
  backend:
    build: ./backend
    depends_on:
      - mysql
      - redis
    
  frontend:
    build: ./frontend
    ports:
      - "80:80"
```

## 安全措施

### 用户认证
- JWT token认证
- 密码加密存储
- 登录限制与保护

### 数据安全
- SQL注入防护
- XSS防护
- CSRF防护
- 敏感数据加密

### API安全
- 限流措施
- 请求签名
- HTTPS传输

## 性能优化

### 前端优化
- 路由懒加载
- 组件按需加载
- 静态资源CDN
- 图片懒加载

### 后端优化
- 数据库索引优化
- 缓存策略
- 异步处理
- 分页加载

## AI写作流程

### 1. 写作前准备
- 选择小说类型
- 设定世界观
- 角色设计
- 大纲规划

### 2. AI辅助写作
- 情节展开建议
- 人物对话生成
- 场景描写优化
- 文风调整

### 3. 内容优化
- 文本润色
- 情节连贯性检查
- 人物性格一致性检查
- 错别字检查

## 开发计划

1. 环境搭建
   - 前端项目初始化
   - 后端项目初始化
   - Docker环境配置
   - 数据库初始化

2. 基础功能开发
   - 用户系统
   - 小说管理
   - 阅读器

3. AI功能集成
   - DeepSeek API对接
   - 提示词系统
   - AI写作助手

4. 优化与部署
   - 性能优化
   - 安全加固
   - 部署上线

## 贡献指南

欢迎提交 Issue 和 Pull Request 来帮助改进项目。

## 许可证

[MIT License](LICENSE)