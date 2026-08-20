package main

import (
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// ---------- auth helpers ----------

func hashPassword(pw string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
	return string(b), err
}

func checkPassword(hash, pw string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(pw)) == nil
}

func jwtSecret() []byte {
	s := os.Getenv("PAYMENT_JWT_SECRET")
	if s == "" {
		s = "dev-secret-change-me"
	}
	return []byte(s)
}

func signToken(userID string) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(72 * time.Hour).Unix(),
		"iat": time.Now().Unix(),
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(jwtSecret())
}

func userFromToken(tok string) (*User, error) {
	t, err := jwt.Parse(tok, func(t *jwt.Token) (any, error) {
		return jwtSecret(), nil
	}, jwt.WithValidMethods([]string{"HS256"}))
	if err != nil || !t.Valid {
		return nil, err
	}
	claims, ok := t.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errNoUser
	}
	sub, _ := claims["sub"].(string)
	var u User
	if err := DB.First(&u, "id = ?", sub).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

var errNoUser = fiber.NewError(fiber.StatusUnauthorized, "unauthorized")

// requireAuth returns the authenticated user via ctx.Locals("user").
func requireAuth(c fiber.Ctx) error {
	h := c.Get("Authorization")
	if len(h) < 8 || h[:7] != "Bearer " {
		return errNoUser
	}
	u, err := userFromToken(h[7:])
	if err != nil {
		return errNoUser
	}
	c.Locals("user", u)
	return c.Next()
}

func me(c fiber.Ctx) *User {
	u, _ := c.Locals("user").(*User)
	return u
}

// ---------- app ----------

func newApp() *fiber.App {
	app := fiber.New(fiber.Config{
		AppName:      "payment-gateway",
		ErrorHandler: errHandler,
	})

	app.Use(logger.New(logger.Config{
		Format: "${time} ${status} ${method} ${path} ${latency}\n",
	}))
	app.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{"Authorization", "Content-Type"},
		AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
	}))

	api := app.Group("/api/v1")

	// public
	api.Get("/health", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{"app": "payment-gateway", "status": "ok"})
	})
	api.Post("/auth/login", handleLogin)

	// protected
	authed := api.Group("", requireAuth)
	authed.Get("/auth/me", handleMe)
	authed.Patch("/auth/password", handleChangePassword)

	authed.Get("/users", handleListUsers)
	authed.Post("/users", handleCreateUser)
	authed.Patch("/users/:id", handleUpdateUser)
	authed.Delete("/users/:id", handleDeleteUser)

	authed.Get("/wallets", handleListWallets)
	authed.Post("/wallets", handleCreateWallet)
	authed.Patch("/wallets/:id", handleUpdateWallet)
	authed.Delete("/wallets/:id", handleDeleteWallet)

	authed.Get("/providers", handleListProviders)
	authed.Post("/providers", handleCreateProvider)
	authed.Patch("/providers/:id", handleUpdateProvider)
	authed.Delete("/providers/:id", handleDeleteProvider)

	mountSPA(app)

	return app
}

func errHandler(c fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	msg := "internal server error"
	if fe, ok := err.(*fiber.Error); ok {
		code = fe.Code
		msg = fe.Message
	} else if err == gorm.ErrRecordNotFound {
		code = fiber.StatusNotFound
		msg = "not found"
	} else {
		log.Printf("error: %v", err)
	}
	return c.Status(code).JSON(fiber.Map{"error": msg})
}

// ---------- auth handlers ----------

func handleLogin(c fiber.Ctx) error {
	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.Bind().Body(&body); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}
	if body.Email == "" || body.Password == "" {
		return fiber.NewError(fiber.StatusBadRequest, "email and password are required")
	}
	var u User
	if err := DB.First(&u, "email = ?", body.Email).Error; err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "invalid credentials")
	}
	if !checkPassword(u.PasswordHash, body.Password) {
		return fiber.NewError(fiber.StatusUnauthorized, "invalid credentials")
	}
	if !u.Active {
		return fiber.NewError(fiber.StatusForbidden, "account is disabled")
	}
	tok, err := signToken(u.ID)
	if err != nil {
		return err
	}
	return c.JSON(fiber.Map{
		"token": tok,
		"user":  u,
	})
}

func handleMe(c fiber.Ctx) error {
	return c.JSON(me(c))
}

func handleChangePassword(c fiber.Ctx) error {
	u := me(c)
	var body struct {
		OldPassword string `json:"oldPassword"`
		NewPassword string `json:"newPassword"`
	}
	if err := c.Bind().Body(&body); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}
	if !checkPassword(u.PasswordHash, body.OldPassword) {
		return fiber.NewError(fiber.StatusUnauthorized, "current password is incorrect")
	}
	if len(body.NewPassword) < 8 {
		return fiber.NewError(fiber.StatusBadRequest, "new password must be at least 8 characters")
	}
	h, err := hashPassword(body.NewPassword)
	if err != nil {
		return err
	}
	u.PasswordHash = h
	if err := DB.Save(u).Error; err != nil {
		return err
	}
	return c.JSON(fiber.Map{"ok": true})
}
