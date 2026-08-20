package main

import (
	"strings"

	"github.com/gofiber/fiber/v3"
)

var validProviderTypes = map[string]bool{
	"midtrans": true, "xendit": true, "tripay": true, "duitku": true,
	"paypal": true, "stripe": true, "other": true,
}

// maskSecret shows only the first 6 and last 4 characters for display.
func maskSecret(s string) string {
	if len(s) <= 12 {
		return "••••••••"
	}
	return s[:6] + "••••••" + s[len(s)-4:]
}

// publicProvider builds the safe JSON shape — credentials never leave the server.
func publicProvider(p *Provider) fiber.Map {
	return fiber.Map{
		"id":          p.ID,
		"name":        p.Name,
		"type":        p.Type,
		"isProduction": p.IsProduction,
		"merchantId":  maskSecret(p.MerchantID),
		"apiKey":      maskSecret(p.ApiKey),
		"apiSecret":   maskSecret(p.ApiSecret),
		"webhookKey":  maskSecret(p.WebhookKey),
		"enabled":     p.Enabled,
		"createdAt":   p.CreatedAt,
		"updatedAt":   p.UpdatedAt,
	}
}

func handleListProviders(c fiber.Ctx) error {
	var providers []Provider
	if err := DB.Order("created_at DESC").Find(&providers).Error; err != nil {
		return err
	}
	out := make([]fiber.Map, 0, len(providers))
	for i := range providers {
		out = append(out, publicProvider(&providers[i]))
	}
	return c.JSON(fiber.Map{"providers": out})
}

func validateProviderType(t string) bool {
	return validProviderTypes[t]
}

func handleCreateProvider(c fiber.Ctx) error {
	var body struct {
		Name         string `json:"name"`
		Type         string `json:"type"`
		IsProduction bool   `json:"isProduction"`
		ApiKey       string `json:"apiKey"`
		ApiSecret    string `json:"apiSecret"`
		MerchantID   string `json:"merchantId"`
		WebhookKey   string `json:"webhookKey"`
		Enabled      *bool  `json:"enabled"`
	}
	if err := c.Bind().Body(&body); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}
	if body.Name == "" {
		return fiber.NewError(fiber.StatusBadRequest, "name is required")
	}
	if body.Type == "" {
		body.Type = "midtrans"
	}
	if !validateProviderType(body.Type) {
		return fiber.NewError(fiber.StatusBadRequest, "unsupported provider type")
	}
	var cnt int64
	DB.Model(&Provider{}).Where("name = ?", body.Name).Count(&cnt)
	if cnt > 0 {
		return fiber.NewError(fiber.StatusConflict, "provider with this name already exists")
	}
	enabled := true
	if body.Enabled != nil {
		enabled = *body.Enabled
	}
	ak, err := encryptSecret(body.ApiKey)
	if err != nil {
		return err
	}
	as, err := encryptSecret(body.ApiSecret)
	if err != nil {
		return err
	}
	mid, err := encryptSecret(body.MerchantID)
	if err != nil {
		return err
	}
	wk, err := encryptSecret(body.WebhookKey)
	if err != nil {
		return err
	}
	p := Provider{
		Name:         body.Name,
		Type:         body.Type,
		IsProduction: body.IsProduction,
		ApiKey:       ak,
		ApiSecret:    as,
		MerchantID:   mid,
		WebhookKey:   wk,
		Enabled:      enabled,
	}
	if err := DB.Create(&p).Error; err != nil {
		return err
	}
	return c.Status(fiber.StatusCreated).JSON(publicProvider(&p))
}

func handleUpdateProvider(c fiber.Ctx) error {
	var p Provider
	if err := DB.First(&p, "id = ?", c.Params("id")).Error; err != nil {
		return fiber.NewError(fiber.StatusNotFound, "provider not found")
	}
	var body struct {
		Name         *string `json:"name"`
		Type         *string `json:"type"`
		IsProduction *bool   `json:"isProduction"`
		ApiKey       *string `json:"apiKey"`
		ApiSecret    *string `json:"apiSecret"`
		MerchantID   *string `json:"merchantId"`
		WebhookKey   *string `json:"webhookKey"`
		Enabled      *bool   `json:"enabled"`
	}
	if err := c.Bind().Body(&body); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}
	if body.Name != nil {
		if strings.TrimSpace(*body.Name) == "" {
			return fiber.NewError(fiber.StatusBadRequest, "name cannot be empty")
		}
		var cnt int64
		DB.Model(&Provider{}).Where("name = ? AND id <> ?", *body.Name, p.ID).Count(&cnt)
		if cnt > 0 {
			return fiber.NewError(fiber.StatusConflict, "provider with this name already exists")
		}
		p.Name = *body.Name
	}
	if body.Type != nil {
		if !validateProviderType(*body.Type) {
			return fiber.NewError(fiber.StatusBadRequest, "unsupported provider type")
		}
		p.Type = *body.Type
	}
	if body.IsProduction != nil {
		p.IsProduction = *body.IsProduction
	}
	if body.Enabled != nil {
		p.Enabled = *body.Enabled
	}
	// encrypt + store any credential that was provided (non-empty)
	if body.ApiKey != nil && *body.ApiKey != "" {
		enc, err := encryptSecret(*body.ApiKey)
		if err != nil {
			return err
		}
		p.ApiKey = enc
	}
	if body.ApiSecret != nil && *body.ApiSecret != "" {
		enc, err := encryptSecret(*body.ApiSecret)
		if err != nil {
			return err
		}
		p.ApiSecret = enc
	}
	if body.MerchantID != nil && *body.MerchantID != "" {
		enc, err := encryptSecret(*body.MerchantID)
		if err != nil {
			return err
		}
		p.MerchantID = enc
	}
	if body.WebhookKey != nil && *body.WebhookKey != "" {
		enc, err := encryptSecret(*body.WebhookKey)
		if err != nil {
			return err
		}
		p.WebhookKey = enc
	}
	if err := DB.Save(&p).Error; err != nil {
		return err
	}
	return c.JSON(publicProvider(&p))
}

func handleDeleteProvider(c fiber.Ctx) error {
	var p Provider
	if err := DB.First(&p, "id = ?", c.Params("id")).Error; err != nil {
		return fiber.NewError(fiber.StatusNotFound, "provider not found")
	}
	if err := DB.Delete(&p).Error; err != nil {
		return err
	}
	return c.JSON(fiber.Map{"ok": true})
}
