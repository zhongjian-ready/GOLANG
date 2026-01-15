package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"

	"github.com/zhongjian-ready/goapi/internal/config"
	"github.com/zhongjian-ready/goapi/internal/database"
	"github.com/zhongjian-ready/goapi/internal/database/mocks"
	"github.com/zhongjian-ready/goapi/internal/handlers"
	"github.com/zhongjian-ready/goapi/internal/models"
)

// 为了测试，我们需要禁用 logrus 的输出，避免测试日志刷屏
func init() {
	log.SetOutput(new(bytes.Buffer))
}

func TestLogin_Success(t *testing.T) {
	// 1. 准备测试数据 (Arrangement)
	// 模拟配置
	cfg := &config.Config{
		JWTSecret: "test_secret_key",
	}

	// 模拟数据库
	mockDB := new(mocks.MockDatabase)

	// 准备密码 hash
	password := "securepassword"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	// 设定 Mock 行为：当调用 GetUserLoginDetails("testuser") 时，返回预设的用户信息
	mockUser := &database.LoginDetails{
		UserID:   1,
		Username: "testuser",
		Password: string(hashedPassword),
	}
	// 使用 mock.Anything 作为 context 参数的占位符，因为 request context 每次都不一样
	mockDB.On("GetUserLoginDetails", mock.Anything, "testuser").Return(mockUser, nil)

	// 初始化被测对象
	h := handlers.NewAPI(mockDB, cfg)

	// 构造 HTTP 请求
	loginPayload := models.LoginParams{
		Username: "testuser",
		Password: password,
	}
	body, _ := json.Marshal(loginPayload)
	req, _ := http.NewRequest("POST", "/account/login", bytes.NewBuffer(body))

	// 构造 ResponseRecorder 来记录响应
	rr := httptest.NewRecorder()

	// 2. 执行测试 (Action)
	// 直接调用 Handler 方法，或者通过 router 调
	// 这里直接调 Handler 方法更纯粹
	http.HandlerFunc(h.Login).ServeHTTP(rr, req)

	// 3. 验证结果 (Assertion)
	// 验证状态码
	assert.Equal(t, http.StatusOK, rr.Code)

	// 验证响应内容
	var resp models.LoginResponse
	err := json.NewDecoder(rr.Body).Decode(&resp)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.Code)
	assert.NotEmpty(t, resp.Token)

	// 验证 Mock 是否按预期被调用
	mockDB.AssertExpectations(t)
}

func TestLogin_InvalidPassword(t *testing.T) {
	// 1. Set up
	cfg := &config.Config{JWTSecret: "test_secret_key"}
	mockDB := new(mocks.MockDatabase)

	password := "correct_password"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	mockUser := &database.LoginDetails{
		UserID:   1,
		Username: "testuser",
		Password: string(hashedPassword),
	}
	mockDB.On("GetUserLoginDetails", mock.Anything, "testuser").Return(mockUser, nil)

	h := handlers.NewAPI(mockDB, cfg)

	// WRONG password
	loginPayload := models.LoginParams{
		Username: "testuser",
		Password: "wrong_password",
	}
	body, _ := json.Marshal(loginPayload)
	req, _ := http.NewRequest("POST", "/account/login", bytes.NewBuffer(body))
	rr := httptest.NewRecorder()

	// 2. Action
	http.HandlerFunc(h.Login).ServeHTTP(rr, req)

	// 3. Assert
	assert.Equal(t, http.StatusUnauthorized, rr.Code)
	mockDB.AssertExpectations(t)
}
