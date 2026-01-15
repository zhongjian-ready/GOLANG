# Token 鉴权流程与前后端交互 (生产环境标准指南)

目前业界最主流的 API 鉴权方案是 **JWT (JSON Web Token)**。

下面梳理了标准的**前后端鉴权交互流程**：

### 1. 核心流程图解

#### 第一阶段：登录 (获取 Token)

1.  **用户操作**：在前端（网页/App）输入 **账号** 和 **密码**，点击登录。
2.  **前端**：将账号密码发送给后端接口（例如 `POST /login`）。
3.  **后端**：
    - 去数据库查询该用户。
    - **比对密码**（注意：生产环境数据库里存的是**加盐哈希**后的乱码，不是明文 `123456`。后端会用算法比对）。
    - 如果验证通过，后端生成一个 **JWT Token**。
      - 这个 Token 里包含用户信息（如 UserID）、过期时间（如 2 小时后过期）和服务器的数字签名。
4.  **后端**：将这个 Token 返回给前端。

#### 第二阶段：使用 (携带 Token)

1.  **前端**：收到 Token 后，将其存储在浏览器的 `localStorage` 或 `Cookie` 中。
2.  **前端**：用户想要查询余额。前端发起请求 `GET /account/coins`。
    - **关键点**：前端会在 HTTP 请求头（Header）中自动带上 Token：
      ```http
      Authorization: Bearer <你的Token字符串>
      ```
3.  **后端 (中间件层)**：
    - 拦截请求，提取 Header 中的 Token。
    - **验证签名**：确认这个 Token 是我服务器签发的，没有被篡改。
    - **检查过期**：确认 Token 是否超时。
    - 如果都通过，解析出 Token 里的 `UserID`，放入请求上下文 (Context)。
4.  **后端 (业务层)**：`handler` 直接从上下文拿到 `UserID`，查询数据库返回余额。

---

### 2. 代码改造方向

若在 Go 项目中实现这套流程，需要做以下改动：

#### 第一步：引入必要的库

需要用到两个核心库：

- **JWT 库**：用于生成和解析 Token (`github.com/golang-jwt/jwt/v5`)
- **加密库**：用于密码哈希 (`golang.org/x/crypto/bcrypt`)

#### 第二步：修改数据库设计

现在的 `users` 表里存放的是 `auth_token`，这在 JWT 模式下是不需要的。我们需要存**密码哈希**。

```sql
-- 生产环境的表结构示例
CREATE TABLE users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL -- 存放加密后的密码，如 $2a$10$N9qo8uLOickgx2ZMRZoMye...
);
```

#### 第三步：新增 `/login` 接口 (LoginHandler)

需要写一个新的 Handler 来处理登录：

```go
// 伪代码示例
func Login(w http.ResponseWriter, r *http.Request) {
    // 1. 解析请求体中的 username 和 password
    // 2. 查数据库，获取 stored_hash
    // 3. 验证密码：bcrypt.CompareHashAndPassword(stored_hash, password)
    // 4. 如果验证通过，生成 JWT：
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "userid": user.ID,
        "exp":    time.Now().Add(time.Hour * 2).Unix(), // 2小时过期
    })
    tokenString, _ := token.SignedString([]byte("你的密钥"))

    // 5. 返回 token给前端
}
```

#### 第四步：改造 Middleware (Authorization)

现在的中间件是去数据库查 Token，改造后**不需要查库**，直接解密 Token：

```go
// 伪代码示例
func Authorization(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // 1. 获取 Header 里的 Authorization: Bearer xxxx
        // 2. 使用 jwt 库解析 xxxx
        // 3. 如果解析失败或过期 -> 返回 401 Unauthorized
        // 4. 解析成功 -> 拿到 userid，放入 r.Context()
        next.ServeHTTP(w, r)
    })
}
```

### 总结

1.  **不再硬编码 Token**：Token 是动态生成的，临时的。
2.  **数据库不存 Token**：数据库只存密码哈希。
3.  **更安全**：即使 Token 被窃取，它也会过期。即使数据库被拖库，黑客也拿不到用户明文密码。
