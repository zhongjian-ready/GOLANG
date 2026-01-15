package handlers

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"

	"github.com/zhongjian-ready/goapi/api"
	"github.com/zhongjian-ready/goapi/internal/tools"
)

// Login 处理用户的登录请求
func Login(w http.ResponseWriter, r *http.Request) {
	var params api.LoginParams
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		api.RequestErrorHandler(w, err)
		return
	}

	// 1. 获取数据库连接
	database, err := tools.NewDatabase()
	if err != nil {
		api.InternalErrorHandler(w)
		return
	}

	// 2. 根据用户名查找用户
	loginDetails := (*database).GetUserLoginDetails(params.Username)
	if loginDetails == nil {
		log.Errorf("User not found: %s", params.Username)
		api.WriteError(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// 3. 验证密码
	err = bcrypt.CompareHashAndPassword([]byte(loginDetails.Password), []byte(params.Password))
	if err != nil {
		log.Errorf("Invalid password for user: %s", params.Username)
		api.WriteError(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// 4. 生成 JWT Token
	tokenString, err := generateToken(loginDetails.UserID, loginDetails.Username)
	if err != nil {
		log.Error("Failed to generate token:", err)
		api.InternalErrorHandler(w)
		return
	}

	// 5. 返回 Token
	resp := api.LoginResponse{
		Code:  http.StatusOK,
		Token: tokenString,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func generateToken(userID int, username string) (string, error) {
	// 定义 Claims (Token 中包含的数据)
	claims := jwt.MapClaims{
		"userid":   userID,
		"username": username,
		"exp":      time.Now().Add(time.Hour * 2).Unix(), // 2小时后过期
	}

	// 创建 Token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 签名密钥 (生产环境应从环境变量获取)
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "default_secret_key"
	}

	// 生成签名后的字符串
	return token.SignedString([]byte(secret))
}
