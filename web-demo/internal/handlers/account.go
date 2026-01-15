package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"

	"github.com/zhongjian-ready/goapi/internal/models"
	"github.com/zhongjian-ready/goapi/internal/tools/response"
)

// Login 处理用户的登录请求
func (a *API) Login(w http.ResponseWriter, r *http.Request) {
	var params models.LoginParams
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		response.RequestErrorHandler(w, err)
		return
	}

	// 1. 输入验证
	if params.Username == "" || params.Password == "" {
		response.WriteError(w, "Username and password are required", http.StatusBadRequest)
		return
	}

	// 2. 根据用户名查找用户
	loginDetails, err := a.DB.GetUserLoginDetails(r.Context(), params.Username)
	if err != nil {
		log.Error("Database error during login:", err)
		response.InternalErrorHandler(w)
		return
	}

	if loginDetails == nil {
		log.Errorf("User not found: %s", params.Username)
		response.WriteError(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// 3. 验证密码
	err = bcrypt.CompareHashAndPassword([]byte(loginDetails.Password), []byte(params.Password))
	if err != nil {
		log.Errorf("Invalid password for user: %s", params.Username)
		response.WriteError(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// 4. 生成 JWT Token
	tokenString, err := generateToken(r.Context(), loginDetails.UserID, loginDetails.Username, a.Config.JWTSecret)
	if err != nil {
		log.Error("Failed to generate token:", err)
		response.InternalErrorHandler(w)
		return
	}

	// 5. 返回 Token
	resp := models.LoginResponse{
		Code:  http.StatusOK,
		Token: tokenString,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func generateToken(ctx context.Context, userID int, username string, secret string) (string, error) {
	// 定义 Claims (Token 中包含的数据)
	claims := jwt.MapClaims{
		"userid":   userID,
		"username": username,
		"exp":      time.Now().Add(time.Hour * 2).Unix(), // 2小时后过期
	}

	// 创建 Token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 生成签名后的字符串
	return token.SignedString([]byte(secret))
}
