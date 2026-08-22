package main

import (
	"encoding/json"
	"strings"

	"github.com/gofiber/fiber/v3"
)

// maskSecret shows only the first 6 and last 4 characters for display.
func maskSecret(s string) string {
	if len(s) <= 12 {
		return "••••••••"
	}
	return s[:6] + "••••••" + s[len(s)-4:]
}

// publicProvider builds the safe JSON shape — credentials never leave the server.
// Returns the method label + a masked map of the credential fields (non-empty only).
func publicProvider(p *Provider) fiber.Map {
	provider, value, _ := resolveMethod(p.Method)
	label := providerMethodLabel(provider, value)

	// parse encrypted creds -> map; decrypt each non-empty value, then mask.
	creds := map[string]string{}
	_ = json.Unmarshal([]byte(p.Creds), &creds)
	masked := fiber.Map{}
	for k, v := range creds {
		if v == "" {
			continue
		}
		if plain, err := decryptSecret(v); err == nil && plain != "" {
			masked[k] = maskSecret(plain)
		} else {
			masked[k] = maskSecret(v)
		}
	}

	return fiber.Map{
		"id":        p.ID,
		"name":      p.Name,
		"method":    p.Method,
		"provider":  provider,
		"label":     label,
		"creds":     masked,
		"enabled":   p.Enabled,
		"createdAt": p.CreatedAt,
		"updatedAt": p.UpdatedAt,
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

func handleCreateProvider(c fiber.Ctx) error {
	var body struct {
		Name    string            `json:"name"`
		Method  string            `json:"method"` // "provider|value"
		Creds   map[string]string `json:"creds"`
		Enabled *bool             `json:"enabled"`
	}
	if err := c.Bind().Body(&body); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}
	if strings.TrimSpace(body.Name) == "" {
		return fiber.NewError(fiber.StatusBadRequest, "name is required")
	}
	provider, _, ok := resolveMethod(body.Method)
	if !ok {
		return fiber.NewError(fiber.StatusBadRequest, "invalid method, pick one from the catalog")
	}
	_ = provider // provider is implicit in the method string; kept for clarity
	var cnt int64
	DB.Model(&Provider{}).Where("name = ?", body.Name).Count(&cnt)
	if cnt > 0 {
		return fiber.NewError(fiber.StatusConflict, "provider with this name already exists")
	}
	enabled := true
	if body.Enabled != nil {
		enabled = *body.Enabled
	}

	// encrypt each credential value before storage
	enc := map[string]string{}
	for k, v := range body.Creds {
		if strings.TrimSpace(v) == "" {
			continue
		}
		e, err := encryptSecret(v)
		if err != nil {
			return err
		}
		enc[k] = e
	}
	encJSON, err := json.Marshal(enc)
	if err != nil {
		return err
	}

	p := Provider{
		Name:    strings.TrimSpace(body.Name),
		Method:  body.Method,
		Creds:   string(encJSON),
		Enabled: enabled,
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
		Name    *string            `json:"name"`
		Method  *string            `json:"method"`
		Creds   map[string]string  `json:"creds"`
		Enabled *bool              `json:"enabled"`
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
		p.Name = strings.TrimSpace(*body.Name)
	}
	if body.Method != nil {
		if _, _, ok := resolveMethod(*body.Method); !ok {
			return fiber.NewError(fiber.StatusBadRequest, "invalid method, pick one from the catalog")
		}
		p.Method = *body.Method
	}
	if body.Enabled != nil {
		p.Enabled = *body.Enabled
	}
	// merge/overwrite creds: only non-empty values are updated; empty = keep existing
	if body.Creds != nil {
		existing := map[string]string{}
		_ = json.Unmarshal([]byte(p.Creds), &existing)
		for k, v := range body.Creds {
			if strings.TrimSpace(v) == "" {
				continue
			}
			e, err := encryptSecret(v)
			if err != nil {
				return err
			}
			existing[k] = e
		}
		encJSON, err := json.Marshal(existing)
		if err != nil {
			return err
		}
		p.Creds = string(encJSON)
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
