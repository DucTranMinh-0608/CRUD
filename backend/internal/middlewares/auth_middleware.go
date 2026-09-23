package middlewares

import (
	"errors"
	"strings"

	"backend/internal/logger"
	"backend/internal/models"
	"backend/internal/services"
	"backend/internal/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const CurrentUserKey = "currentUser"

func ExtractToken(c *gin.Context) string {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return ""
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) == 2 && strings.EqualFold(parts[0], "Bearer") {
		return strings.TrimSpace(parts[1])
	}
	if len(parts) == 1 {
		return strings.TrimSpace(parts[0])
	}
	return ""
}

func AuthMiddleware(authService services.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := ExtractToken(c)

		if token == "" {
			if logger.Log != nil {
				logger.Log.Warn("Middleware chặn request: thiếu token xác thực", zap.String("path", c.Request.URL.Path))
			}
			utils.UnauthorizedResponse(c, services.ErrMissingToken)
			c.Abort()
			return
		}

		user, err := authService.GetMe(c.Request.Context(), token)
		if err != nil {
			if errors.Is(err, services.ErrInvalidToken) || errors.Is(err, services.ErrMissingToken) || errors.Is(err, services.ErrAccountDisabled) {
				if logger.Log != nil {
					logger.Log.Warn("Middleware chặn request: token không hợp lệ hoặc tài khoản vô hiệu", zap.Error(err), zap.String("path", c.Request.URL.Path))
				}
				utils.UnauthorizedResponse(c, err)
				c.Abort()
				return
			}

			if errors.Is(err, services.ErrUserNotFound) {
				if logger.Log != nil {
					logger.Log.Warn("Middleware chặn request: không tìm thấy người dùng", zap.Error(err), zap.String("path", c.Request.URL.Path))
				}
				utils.NotFoundDataResponse(c, err)
				c.Abort()
				return
			}

			if logger.Log != nil {
				logger.Log.Error("Middleware lỗi hệ thống khi xác thực token", zap.Error(err), zap.String("path", c.Request.URL.Path))
			}
			utils.NoConnectDataResponse(c, err)
			c.Abort()
			return
		}

		c.Set(CurrentUserKey, user)
		c.Next()
	}
}

func RequireRoles(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		val, exists := c.Get(CurrentUserKey)
		if !exists {
			if logger.Log != nil {
				logger.Log.Warn("Role middleware chặn: không tìm thấy thông tin user trong context")
			}
			utils.UnauthorizedResponse(c, services.ErrMissingToken)
			c.Abort()
			return
		}

		user, ok := val.(*models.User)
		if !ok || user == nil {
			if logger.Log != nil {
				logger.Log.Warn("Role middleware chặn: thông tin user không hợp lệ")
			}
			utils.UnauthorizedResponse(c, services.ErrMissingToken)
			c.Abort()
			return
		}

		userRole := strings.ToLower(strings.TrimSpace(user.ChucVu))
		hasPermission := false
		for _, role := range allowedRoles {
			if strings.EqualFold(userRole, strings.TrimSpace(role)) {
				hasPermission = true
				break
			}
		}

		if !hasPermission {
			if logger.Log != nil {
				logger.Log.Warn(
					"Người dùng không đủ quyền hạn truy cập API",
					zap.Int64("user_id", user.ID),
					zap.String("user_role", user.ChucVu),
					zap.Strings("required_roles", allowedRoles),
					zap.String("path", c.Request.URL.Path),
				)
			}
			utils.ForbiddenResponse(c, services.ErrForbidden)
			c.Abort()
			return
		}

		c.Next()
	}
}

func GetCurrentUser(c *gin.Context) (*models.User, bool) {
	val, exists := c.Get(CurrentUserKey)
	if !exists {
		return nil, false
	}
	user, ok := val.(*models.User)
	return user, ok
}
