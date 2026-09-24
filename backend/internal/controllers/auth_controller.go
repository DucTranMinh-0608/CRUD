package controllers

import (
	"errors"

	"go.uber.org/zap"

	"backend/internal/dto"
	"backend/internal/logger"
	"backend/internal/mappers"
	"backend/internal/middlewares"
	"backend/internal/services"
	"backend/internal/utils"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	service services.AuthService
}

func NewAuthController(service services.AuthService) *AuthController {
	return &AuthController{
		service: service,
	}
}

func (ctrl *AuthController) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Log.Warn(
			"Dữ liệu đăng ký không hợp lệ",
			zap.Error(err),
		)
		utils.InvalidDataResponse(c, err)
		return
	}

	registerModel := mappers.ToRegister(&req)

	user, err := ctrl.service.Register(c.Request.Context(), registerModel)
	if err != nil {
		if errors.Is(err, utils.ErrEmpty) ||
			errors.Is(err, utils.ErrInvalidRole) ||
			errors.Is(err, utils.ErrPasswordTooShort) ||
			errors.Is(err, utils.ErrInvalidEmail) ||
			errors.Is(err, utils.ErrEmailAlreadyExists) {

			logger.Log.Warn(
				"Dữ liệu đăng ký không hợp lệ hoặc đã tồn tại",
				zap.Error(err),
			)
			utils.InvalidDataResponse(c, err)
			return
		}

		logger.Log.Error(
			"Lỗi hệ thống khi đăng ký người dùng",
			zap.Error(err),
		)
		utils.NoConnectDataResponse(c, err)
		return
	}

	logger.Log.Info(
		"Đăng ký người dùng thành công",
		zap.String("email", user.Email),
		zap.Int64("user_id", user.ID),
	)

	utils.SuccessResponse(c, "Đăng ký người dùng", mappers.ToUserResponse(user))
}

func (ctrl *AuthController) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Log.Warn(
			"Dữ liệu đăng nhập không hợp lệ",
			zap.Error(err),
		)
		utils.InvalidDataResponse(c, err)
		return
	}

	loginModel := mappers.ToLogin(&req)

	user, token, err := ctrl.service.Login(c.Request.Context(), loginModel)
	if err != nil {
		if errors.Is(err, utils.ErrEmpty) || errors.Is(err, utils.ErrInvalidEmail) {
			logger.Log.Warn(
				"Dữ liệu đăng nhập không hợp lệ",
				zap.Error(err),
			)
			utils.InvalidDataResponse(c, err)
			return
		}

		if errors.Is(err, utils.ErrInvalidCredentials) {
			logger.Log.Warn(
				"Đăng nhập thất bại: sai email hoặc mật khẩu",
				zap.String("email", req.Email),
			)
			utils.UnauthorizedResponse(c, err)
			return
		}

		if errors.Is(err, utils.ErrAccountDisabled) {
			logger.Log.Warn(
				"Đăng nhập thất bại: tài khoản đã bị vô hiệu hóa",
				zap.String("email", req.Email),
			)
			utils.UnauthorizedResponse(c, err)
			return
		}

		if errors.Is(err, utils.ErrUserNotFound) {
			logger.Log.Warn(
				"Không tìm thấy người dùng trong database",
				zap.String("email", req.Email),
			)
			utils.NotFoundDataResponse(c, err)
			return
		}

		logger.Log.Error(
			"Lỗi hệ thống khi đăng nhập",
			zap.Error(err),
		)
		utils.NoConnectDataResponse(c, err)
		return
	}

	logger.Log.Info(
		"Đăng nhập thành công",
		zap.String("email", user.Email),
		zap.Int64("user_id", user.ID),
	)

	utils.SuccessResponse(c, "Đăng nhập", gin.H{
		"access_token": token,
		"user":         mappers.ToUserResponse(user),
	})
}

func (ctrl *AuthController) GetMe(c *gin.Context) {
	user, ok := middlewares.GetCurrentUser(c)
	if !ok || user == nil {
		logger.Log.Warn("Không tìm thấy thông tin người dùng trong context (yêu cầu AuthMiddleware)")
		utils.UnauthorizedResponse(c, utils.ErrMissingToken)
		return
	}

	utils.SuccessResponse(c, "Lấy thông tin người dùng", mappers.ToUserResponse(user))
}

func (ctrl *AuthController) Logout(c *gin.Context) {
	token := middlewares.ExtractToken(c)
	if token == "" {
		if logger.Log != nil {
			logger.Log.Warn("Đăng xuất thất bại: không tìm thấy token xác thực")
		}
		utils.UnauthorizedResponse(c, utils.ErrMissingToken)
		return
	}

	if err := ctrl.service.Logout(c.Request.Context(), token); err != nil {
		if logger.Log != nil {
			logger.Log.Warn(
				"Lỗi khi xoá access token khỏi database",
				zap.Error(err),
			)
		}
		utils.NoConnectDataResponse(c, err)
		return
	}

	if logger.Log != nil {
		logger.Log.Info("Đăng xuất thành công")
	}
	utils.SuccessResponse(c, "Đăng xuất", nil)
}
