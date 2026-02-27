# Go Web 框架

这是一个基于 `net/http` 构建的轻量级 Go Web 框架。

## 功能规划

- [x] **核心功能**
    - [x] 基于 `net/http` 封装核心引擎
    - [x] 提供统一的日志记录功能

- [X] **上下文（Context）**
    - [X] 提供对请求参数、路径参数、表单数据和 JSON 数据的便捷访问方法
    - [X] 提供快速生成 JSON、HTML、字符串等响应的方法

- [X] **路由（Routing）**
    - [X] 支持静态路由（如 `/about`, `/contact`）
    - [X] 支持动态路由/参数路由（如 `/users/:id`）
    - [X] 支持路由分组，可以对一组路由应用相同的前缀或中间件
    - [X] 支持 RESTful 风格的路由（GET, POST, PUT, DELETE 等）

- [X] **中间件（Middleware）**
    - [X] 支持全局中间件
    - [X] 支持单个路由的中间件
    - [X] 支持路由分组的中间件
    - [X] 内置常用的中间件（如 Logger, Recovery）

- [X] **错误处理**
    - [X] 提供统一的 panic 恢复和错误处理机制

