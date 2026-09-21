package controllers

import (
	"errors"
	"strconv"

	"go.uber.org/zap"

	"backend/internal/dto"
	"backend/internal/logger"
	"backend/internal/mappers"
	"backend/internal/services"
	"backend/internal/utils"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	service services.UserService
}

func NewUserController(service services.UserService) *UserController {
	return &UserController{
		service: service,
	}
}

func (ctrl *UserController) GetAllUsers(c *gin.Context) {
	users, err := ctrl.service.GetAllUsers(c.Request.Context())
	if err != nil {

		logger.Log.Error(
			"Lỗi hệ thống khi lấy danh sách người dùng",
			zap.Error(err),
		)

		utils.NoConnectDataResponse(c, err)
		return
	}

	logger.Log.Info(
		"Lấy danh sách người dùng thành công",
		zap.Any("count", len(users)),
	)

	responses := make([]dto.UserResponse, 0, len(users))
	for _, p := range users {
		responses = append(responses, *mappers.ToUserResponse(&p))
	}
	utils.SuccessResponse(c, "Lấy người dùng", responses)
}

func (ctrl *UserController) GetUserByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {

		logger.Log.Warn(
			"ID người dùng không hợp lệ",
			zap.Error(err),
		)

		utils.InvalidDataResponse(c, err)
		return
	}

	user, err := ctrl.service.GetUserByID(c.Request.Context(), id)
	if err != nil {

		if errors.Is(err, services.ErrUserNotFound) {

			logger.Log.Warn(
				"Không tìm thấy người dùng với ID đã cung cấp",
				zap.Int64("id", id),
				zap.Error(err),
			)

			utils.NotFoundDataResponse(c, err)
			return
		}

		logger.Log.Error(
			"Lỗi hệ thống khi tìm kiếm người dùng",
			zap.Error(err),
		)

		utils.NoConnectDataResponse(c, err)
		return
	}

	logger.Log.Info(
		"Tìm thấy người dùng với ID đã cung cấp",
		zap.Int64("id", id),
	)

	utils.SuccessResponse(c, "Lấy người dùng", mappers.ToUserResponse(user))
}

func (ctrl *UserController) UpdateUser(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {

		logger.Log.Warn(
			"ID người dùng không hợp lệ",
			zap.Error(err),
		)

		utils.InvalidDataResponse(c, err)
		return
	}

	var req dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {

		logger.Log.Warn(
			"Dữ liệu cập nhật không hợp lệ",
			zap.Error(err),
		)

		utils.InvalidDataResponse(c, err)
		return
	}

	updatedUserModel := mappers.ToUpdateUser(&req)

	updatedUser, err := ctrl.service.UpdateUser(c.Request.Context(), id, updatedUserModel)
	if err != nil {
		if errors.Is(err, services.ErrUserNotFound) {

			logger.Log.Warn(
				"Không tìm thấy người dùng để cập nhật",
				zap.Int64("id", id),
				zap.Error(err),
			)

			utils.NotFoundDataResponse(c, err)
			return
		}

		if errors.Is(err, services.ErrInvalidRole) ||
			errors.Is(err, services.ErrInvalidStatus) ||
			errors.Is(err, services.ErrEmpty) {

			logger.Log.Warn(
				"Dữ liệu cập nhật người dùng không hợp lệ",
				zap.Error(err),
			)

			utils.InvalidDataResponse(c, err)
			return
		}

		logger.Log.Error(
			"Lỗi hệ thống khi cập nhật người dùng",
			zap.Error(err),
		)

		utils.NoConnectDataResponse(c, err)
		return
	}

	logger.Log.Info(
		"Cập nhật người dùng thành công",
		zap.Int64("id", id),
	)

	utils.SuccessResponse(c, "Cập nhật người dùng", mappers.ToUserResponse(updatedUser))
}
