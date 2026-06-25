package api

import (
	"net/http"

	"catering-backend/internal/model"
	"catering-backend/pkg/database"
	"catering-backend/pkg/util"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// UserRequest 用于接收前端传来的注册和登录 JSON 参数
type UserRequest struct {
	Username string `json:"username" binding:"required,min=3,max=20"`
	Password string `json:"password" binding:"required,min=6"`
}

// Register 用户注册接口
func Register(c *gin.Context) {
	var req UserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "参数格式错误，账号最少3位，密码最少6位"})
		return
	}

	// 1. 检查数据库中用户名是否已存在
	var count int64
	database.DB.Model(&model.User{}).Where("username = ?", req.Username).Count(&count)
	if count > 0 {
		c.JSON(http.StatusOK, gin.H{"code": 409, "msg": "用户名已被注册，请换一个重试"})
		return
	}

	// 2. 对密码进行 Bcrypt 不可逆加密
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "系统内部错误，密码加密失败"})
		return
	}

	// 3. 构造模型并存入 MySQL 数据库
	newUser := model.User{
		Username:     req.Username,
		PasswordHash: string(hashedPassword),
	}
	if err := database.DB.Create(&newUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "注册失败，写入数据库异常"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "注册成功！"})
}

// Login 用户登录接口
func Login(c *gin.Context) {
	var req UserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "参数缺失"})
		return
	}

	// 1. 根据用户名从数据库获取用户记录
	var user model.User
	if err := database.DB.Where("username = ?", req.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 404, "msg": "用户不存在，请先注册"})
		return
	}

	// 2. 校验前端传来的明文密码与数据库里的 Hash 密文是否匹配
	err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 401, "msg": "密码错误，请重试"})
		return
	}

	// 3. 密码校验成功，生成 JWT Token 令牌
	token, err := util.GenerateToken(user.ID, user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "登录状态签发失败"})
		return
	}

	// 4. 下发 Token 与基础用户信息
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "登录成功",
		"data": gin.H{
			"token":    token,
			"username": user.Username,
			"id":       user.ID,
			"role":     user.Role,
		},
	})
}
