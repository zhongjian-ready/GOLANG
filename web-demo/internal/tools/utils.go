package tools

import "os"

// GetJWTSecret 获取 JWT 签名密钥
// 优先从环境变量 JWT_SECRET 获取，如果未设置则使用默认值。
// 注意：在生产环境中，必须确保环境变量已设置，严禁使用默认值。
func GetJWTSecret() []byte {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "default_secret_key"
	}
	return []byte(secret)
}
