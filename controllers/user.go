package controllers

import (
	"fmt"
	"strings"

	"github.com/abe27/api/crypto/configs"
	"github.com/abe27/api/crypto/models"
	"github.com/abe27/api/crypto/services"
	"github.com/gofiber/fiber/v2"
)

func ProfileController(c *fiber.Ctx) error {
	var r models.Response
	s := c.Get("Authorization")
	token := strings.TrimPrefix(s, "Bearer ")
	if token == "" {
		r.Message = "Token required"
		CreateLogger("ร้องขอดูข้อมูลส่วนบุคคล", r.Message)
		return c.Status(fiber.StatusUnauthorized).JSON(&r)
	}

	id, er := services.ValidateToken(token)
	if er != nil {
		r.Message = er.Error()
		CreateLogger("ร้องขอดูข้อมูลส่วนบุคคล", r.Message)
		return c.Status(fiber.StatusUnauthorized).JSON(&r)
	}

	uid := fmt.Sprintf("%s", id)
	var user models.User
	if err := configs.Store.Where("id", uid).Find(&user).Error; err != nil {
		r.Message = err.Error()
		CreateLogger("ร้องขอดูข้อมูลส่วนบุคคล", fmt.Sprintf("%s พยายามร้องขอดูข้อมูลส่วนบุคคล", uid))
		return c.Status(fiber.StatusUnauthorized).JSON(&r)
	}

	if user.ID == "" {
		r.Message = "User not found"
		CreateLogger("ร้องขอดูข้อมูลส่วนบุคคล", fmt.Sprintf("%s ไม่พบข้อมูลส่วนบุคคล", uid))
		return c.Status(fiber.StatusUnauthorized).JSON(&r)
	}

	r.Success = true
	r.Data = &user
	CreateLogger("ร้องขอดูข้อมูลส่วนบุคคล", fmt.Sprintf("%s แสดงข้อมูลส่วนบุคล", uid))
	return c.Status(fiber.StatusOK).JSON(&r)
}

func RegisterController(c *fiber.Ctx) error {
	var r models.Response

	var frmUser *models.User
	if err := c.BodyParser(&frmUser); err != nil {
		CreateLogger("ลงทะเบียน", fmt.Sprintf("%s ลงทะเบียนไม่สำเร็จ %s", frmUser.Email, r.Message))
		return c.Status(fiber.StatusBadRequest).JSON(r)
	}

	var user models.User
	if err := configs.Store.Create(&user).Error; err != nil {
		r.Message = err.Error()
		CreateLogger("ลงทะเบียน", fmt.Sprintf("%s ลงทะเบียนไม่สำเร็จ %s", frmUser.Email, r.Message))
		return c.Status(fiber.StatusBadRequest).JSON(r)
	}

	r.Success = true
	r.Data = &user
	// Create Logger
	CreateLogger("ลงทะเบียน", fmt.Sprintf("%s ลงทะเบียนเรียบร้อยแล้ว", frmUser.Email))
	return c.Status(fiber.StatusCreated).JSON(&r)
}

func LoginController(c *fiber.Ctx) error {
	var r models.Response
	var user models.UserLoginForm
	if err := c.BodyParser(&user); err != nil {
		r.Message = err.Error()
		return c.Status(fiber.StatusInternalServerError).JSON(&r)
	}
	// Check AuthorizationRequired
	db := configs.Store
	var userData models.User
	if err := db.Where("email=?", user.Email).First(&userData).Error; err != nil {
		r.Message = err.Error()
		// Create Logger
		CreateLogger("เข้าสู่ระบบ", fmt.Sprintf("%s เกิดข้อผิดพลาด %s", user.Email, err.Error()))
		return c.Status(fiber.StatusNotFound).JSON(&r)
	}

	isMatched := services.CheckPasswordHashing(user.Password, userData.Password)
	if !isMatched {
		r.Message = "Password not match!"
		// Create Logger
		CreateLogger("เข้าสู่ระบบ", fmt.Sprintf("%s Password not match!", user.Email))
		return c.Status(fiber.StatusUnauthorized).JSON(&r)
	}

	// Create Token
	auth, _ := services.CreateToken(userData)
	r.Message = "Auth success!"
	r.Data = &auth
	// Create Logger
	CreateLogger("เข้าสู่ระบบ", fmt.Sprintf("%s เข้าสู่ระบบเรียบร้อย", user.Email))
	return c.Status(fiber.StatusOK).JSON(&r)
}

func LogOutController(c *fiber.Ctx) error {
	var r models.Response
	s := c.Get("Authorization")
	token := strings.TrimPrefix(s, "Bearer ")
	if token == "" {
		r.Message = "Token required"
		CreateLogger("ร้องขอดูข้อมูลส่วนบุคคล", r.Message)
		return c.Status(fiber.StatusUnauthorized).JSON(&r)
	}

	id, er := services.ValidateToken(token)
	if er != nil {
		r.Message = er.Error()
		CreateLogger("ออกจากระบบ", r.Message)
		return c.Status(fiber.StatusUnauthorized).JSON(&r)
	}

	uid := fmt.Sprintf("%s", id)
	// Destry session
	configs.Store.Delete(&models.UserLogin{}, &models.UserLogin{UserID: uid})
	// Create Logger
	r.Success = true
	r.Message = "Logout success!"
	CreateLogger("ออกจากระบบ", fmt.Sprintf("%s ออกจากระบบเรียบร้อย", uid))
	return c.Status(fiber.StatusOK).JSON(&r)
}
