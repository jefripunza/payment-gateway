package main

import (
	"log"

	"github.com/gofiber/fiber/v3"
	"golang.org/x/crypto/bcrypt"
)

func seedAdmin() {
	var count int64
	if err := DB.Model(&User{}).Count(&count).Error; err != nil {
		// Tables may not exist yet when DATABASE_MIGRATE=false on a fresh DB.
		log.Printf("seed admin: skipped (%v) — run once with DATABASE_MIGRATE=true to create schema", err)
		return
	}
	if count > 0 {
		return
	}
	email := "admin@sawang.tech"
	pw := "Admin@123"
	h, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("seed admin: %v", err)
	}
	u := User{Name: "Administrator", Email: email, PasswordHash: string(h), Role: "admin", Active: true}
	if err := DB.Create(&u).Error; err != nil {
		log.Fatalf("seed admin: %v", err)
	}
	log.Printf("seeded admin user: %s", email)
}

// ---------- users ----------

func handleListUsers(c fiber.Ctx) error {
	var users []User
	if err := DB.Order("created_at DESC").Find(&users).Error; err != nil {
		return err
	}
	return c.JSON(fiber.Map{"users": users})
}

func handleCreateUser(c fiber.Ctx) error {
	var body struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
		Role     string `json:"role"`
		Active   *bool  `json:"active"`
	}
	if err := c.Bind().Body(&body); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}
	if body.Name == "" || body.Email == "" {
		return fiber.NewError(fiber.StatusBadRequest, "name and email are required")
	}
	if len(body.Password) < 8 {
		return fiber.NewError(fiber.StatusBadRequest, "password must be at least 8 characters")
	}
	var cnt int64
	DB.Model(&User{}).Where("email = ?", body.Email).Count(&cnt)
	if cnt > 0 {
		return fiber.NewError(fiber.StatusConflict, "user with this email already exists")
	}
	h, err := hashPassword(body.Password)
	if err != nil {
		return err
	}
	role := body.Role
	if role != "admin" && role != "viewer" {
		role = "viewer"
	}
	active := true
	if body.Active != nil {
		active = *body.Active
	}
	u := User{Name: body.Name, Email: body.Email, PasswordHash: h, Role: role, Active: active}
	if err := DB.Create(&u).Error; err != nil {
		return err
	}
	return c.Status(fiber.StatusCreated).JSON(u)
}

func handleUpdateUser(c fiber.Ctx) error {
	var u User
	if err := DB.First(&u, "id = ?", c.Params("id")).Error; err != nil {
		return fiber.NewError(fiber.StatusNotFound, "user not found")
	}
	var body struct {
		Name     *string `json:"name"`
		Email    *string `json:"email"`
		Password *string `json:"password"`
		Role     *string `json:"role"`
		Active   *bool   `json:"active"`
	}
	if err := c.Bind().Body(&body); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}
	if body.Name != nil {
		if *body.Name == "" {
			return fiber.NewError(fiber.StatusBadRequest, "name cannot be empty")
		}
		u.Name = *body.Name
	}
	if body.Email != nil {
		if *body.Email == "" {
			return fiber.NewError(fiber.StatusBadRequest, "email cannot be empty")
		}
		var cnt int64
		DB.Model(&User{}).Where("email = ? AND id <> ?", *body.Email, u.ID).Count(&cnt)
		if cnt > 0 {
			return fiber.NewError(fiber.StatusConflict, "user with this email already exists")
		}
		u.Email = *body.Email
	}
	if body.Password != nil && *body.Password != "" {
		if len(*body.Password) < 8 {
			return fiber.NewError(fiber.StatusBadRequest, "password must be at least 8 characters")
		}
		h, err := hashPassword(*body.Password)
		if err != nil {
			return err
		}
		u.PasswordHash = h
	}
	if body.Role != nil {
		if *body.Role != "admin" && *body.Role != "viewer" {
			return fiber.NewError(fiber.StatusBadRequest, "role must be admin or viewer")
		}
		u.Role = *body.Role
	}
	if body.Active != nil {
		u.Active = *body.Active
	}
	if err := DB.Save(&u).Error; err != nil {
		return err
	}
	return c.JSON(u)
}

func handleDeleteUser(c fiber.Ctx) error {
	cur := me(c)
	if cur != nil && cur.ID == c.Params("id") {
		return fiber.NewError(fiber.StatusBadRequest, "you cannot delete your own account")
	}
	var u User
	if err := DB.First(&u, "id = ?", c.Params("id")).Error; err != nil {
		return fiber.NewError(fiber.StatusNotFound, "user not found")
	}
	if err := DB.Delete(&u).Error; err != nil {
		return err
	}
	return c.JSON(fiber.Map{"ok": true})
}

// ---------- wallets ----------

func handleListWallets(c fiber.Ctx) error {
	var wallets []Wallet
	if err := DB.Order("created_at DESC").Find(&wallets).Error; err != nil {
		return err
	}
	return c.JSON(fiber.Map{"wallets": wallets})
}

func handleCreateWallet(c fiber.Ctx) error {
	var body struct {
		Name     string `json:"name"`
		Currency string `json:"currency"`
		Balance  int64  `json:"balance"`
	}
	if err := c.Bind().Body(&body); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}
	if body.Name == "" {
		return fiber.NewError(fiber.StatusBadRequest, "name is required")
	}
	cur := body.Currency
	if cur == "" {
		cur = "IDR"
	}
	if body.Balance < 0 {
		return fiber.NewError(fiber.StatusBadRequest, "balance cannot be negative")
	}
	w := Wallet{Name: body.Name, Currency: cur, Balance: body.Balance}
	if err := DB.Create(&w).Error; err != nil {
		return err
	}
	return c.Status(fiber.StatusCreated).JSON(w)
}

func handleUpdateWallet(c fiber.Ctx) error {
	var w Wallet
	if err := DB.First(&w, "id = ?", c.Params("id")).Error; err != nil {
		return fiber.NewError(fiber.StatusNotFound, "wallet not found")
	}
	var body struct {
		Name     *string `json:"name"`
		Currency *string `json:"currency"`
		Balance  *int64  `json:"balance"`
	}
	if err := c.Bind().Body(&body); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}
	if body.Name != nil {
		if *body.Name == "" {
			return fiber.NewError(fiber.StatusBadRequest, "name cannot be empty")
		}
		w.Name = *body.Name
	}
	if body.Currency != nil {
		if *body.Currency == "" {
			return fiber.NewError(fiber.StatusBadRequest, "currency cannot be empty")
		}
		w.Currency = *body.Currency
	}
	if body.Balance != nil {
		if *body.Balance < 0 {
			return fiber.NewError(fiber.StatusBadRequest, "balance cannot be negative")
		}
		w.Balance = *body.Balance
	}
	if err := DB.Save(&w).Error; err != nil {
		return err
	}
	return c.JSON(w)
}

func handleDeleteWallet(c fiber.Ctx) error {
	var w Wallet
	if err := DB.First(&w, "id = ?", c.Params("id")).Error; err != nil {
		return fiber.NewError(fiber.StatusNotFound, "wallet not found")
	}
	if err := DB.Delete(&w).Error; err != nil {
		return err
	}
	return c.JSON(fiber.Map{"ok": true})
}
