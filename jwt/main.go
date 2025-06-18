package main

import (
	"fmt"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

// 在之前的一些web项目中，我们通常使用的是Cookie-Session模式实现用户认证。相关流程大致如下：
// 1. 用户在浏览器端填写用户名和密码，并发送给服务端
// 2. 服务端对用户名和密码校验通过后会生成一份保存当前用户相关信息的session数据和一个与之对应的标识（通常称为session_id）
// 3. 服务端返回响应时将上一步的session_id写入用户浏览器的Cookie
// 4. 后续用户来自该浏览器的每次请求都会自动携带包含session_id的Cookie
// 5. 服务端通过请求中的session_id就能找到之前保存的该用户那份session数据，从而获取该用户的相关信息。
// 这种方案依赖于客户端（浏览器）保存Cookie，并且需要在服务端存储用户的session数据。
// 在移动互联网时代，我们的用户可能使用浏览器也可能使用APP来访问我们的服务，我们的web应用可能是前后端分开部署在不同的端口，有时候我们还需要支持第三方登录，这下Cookie-Session的模式就有些力不从心了。
// JWT就是一种基于Token的轻量级认证模式，服务端认证通过后，会生成一个JSON对象，经过签名后得到一个Token（令牌）再发回给用户，用户后续请求只需要带上这个Token，服务端解密之后就能获取该用户的相关信息了。

type Claims struct {
	UserID uint   `json:"user_id"`
	Role   string `json:"role"`
	jwt.StandardClaims
}

// 密钥（生产环境应从安全存储获取）
var secretKey = []byte("your-secret-key")

func GenerateToken(userID uint, role string) (string, error) {
	claims := &Claims{
		UserID: userID,
		Role:   role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(2 * time.Minute).Unix(),
			Issuer:    "your-app",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}

func ParseToken(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secretKey, nil
	})
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, err
}

func main() {
	// 生成token
	token, err := GenerateToken(1, "admin")
	if err != nil {
		panic(err)
	}
	fmt.Println("Generated Token:", token)

	// 解析token
	claims, err := ParseToken(token)
	if err != nil {
		panic(err)
	}
	fmt.Printf("UserID: %d, Role: %s, ExpiresAt: %s\n",
		claims.UserID, claims.Role, time.Unix(claims.ExpiresAt, 0))

	// HTTP服务器示例
	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		token, _ := GenerateToken(1, "user")
		w.Write([]byte(token))
	})

	http.HandleFunc("/protected", func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header missing", http.StatusUnauthorized)
			return
		}
		tokenStr := authHeader[7:] // 移除"Bearer "前缀
		claims, err := ParseToken(tokenStr)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}
		w.Write([]byte(fmt.Sprintf("Welcome, %s! UserID: %d", claims.Role, claims.UserID)))
	})

	http.ListenAndServe(":8080", nil)
}
