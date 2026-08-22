package main

import (
	"github.com/gofiber/fiber/v3"
	"github.com/jefripunza/go-payment-method/tripay"
	"github.com/jefripunza/go-payment-method/xendit"
)

// ProviderMethod describes one selectable payment method for a provider.
type ProviderMethod struct {
	Value string `json:"value"` // raw constant from the library (e.g. "BCA", "OVO", "QRIS")
	Label string `json:"label"` // human-friendly name shown in the dropdown
}

// MethodField describes one dynamic credential input for a provider method.
type MethodField struct {
	Key         string `json:"key"`
	Label       string `json:"label"`
	Required    bool   `json:"required"`
	Placeholder string `json:"placeholder,omitempty"`
}

// CatalogMethod is one selectable entry in the FE methods dropdown.
// Value is the raw library constant (e.g. "QRIS", "CREDIT_CARD");
// Provider tells the backend which gateway it belongs to.
type CatalogMethod struct {
	Provider string       `json:"provider"`
	Value    string       `json:"value"`
	Label    string       `json:"label"`
	Fields   []MethodField `json:"fields"`
}

// providerMethodFields defines the dynamic credential inputs for each gateway.
// All methods of a gateway share its credential set (that's how the
// gateways work — one keypair powers every channel).
var providerMethodFields = map[string][]MethodField{
	"midtrans": {
		{Key: "serverKey", Label: "Server Key", Required: true, Placeholder: "Midtrans server key"},
		{Key: "clientKey", Label: "Client Key", Required: true, Placeholder: "Midtrans client key"},
		{Key: "merchantId", Label: "Merchant ID", Required: false, Placeholder: "Optional"},
	},
	"xendit": {
		{Key: "secretKey", Label: "Secret Key (x-api-key)", Required: true, Placeholder: "Xendit API key"},
		{Key: "publishableKey", Label: "Publishable Key", Required: false, Placeholder: "Optional"},
	},
	"tripay": {
		{Key: "apiKey", Label: "API Key", Required: true, Placeholder: "TriPay API key"},
		{Key: "privateKey", Label: "Private Key", Required: true, Placeholder: "TriPay private key"},
		{Key: "merchantCode", Label: "Merchant Code", Required: true, Placeholder: "TriPay merchant code"},
	},
	"duitku": {
		{Key: "merchantCode", Label: "Merchant Code", Required: true, Placeholder: "Duitku merchant code"},
		{Key: "apiKey", Label: "API Key", Required: true, Placeholder: "Duitku API key"},
	},
	"stripe": {
		{Key: "secretKey", Label: "Secret Key", Required: true, Placeholder: "sk_live_..."},
		{Key: "publishableKey", Label: "Publishable Key", Required: false, Placeholder: "pk_live_... (optional)"},
	},
	"paypal": {
		{Key: "clientId", Label: "Client ID", Required: true, Placeholder: "PayPal client ID"},
		{Key: "clientSecret", Label: "Client Secret", Required: true, Placeholder: "PayPal client secret"},
	},
	"other": {
		{Key: "apiKey", Label: "API Key", Required: true, Placeholder: "API key"},
		{Key: "apiSecret", Label: "API Secret", Required: false, Placeholder: "Optional"},
	},
}

// providerMethods is the single source of truth for the FE dropdown.
// Xendit & Tripay lists come straight from the library constants
// (xendit.PaymentMethod*, tripay.Channel*). Midtrans/Stripe/Duitku have no
// method enum in the library (they accept free-form strings), so we ship a
// sensible manual list of the common Indonesian channels — matching what the
// library's own examples use.
var providerMethods = map[string][]ProviderMethod{
	"midtrans": {
		{Value: "credit_card", Label: "Credit Card"},
		{Value: "bank_transfer", Label: "Bank Transfer"},
		{Value: "bca_va", Label: "BCA Virtual Account"},
		{Value: "bni_va", Label: "BNI Virtual Account"},
		{Value: "bri_va", Label: "BRI Virtual Account"},
		{Value: "mandiri_va", Label: "Mandiri Virtual Account"},
		{Value: "permata_va", Label: "Permata Virtual Account"},
		{Value: "other_va", Label: "Other Virtual Account"},
		{Value: "echannel", Label: "E-Channel (Mandiri)"},
		{Value: "qris", Label: "QRIS"},
		{Value: "gopay", Label: "GoPay"},
		{Value: "shopeepay", Label: "ShopeePay"},
		{Value: "dana", Label: "DANA"},
		{Value: "ovo", Label: "OVO"},
		{Value: "cimb_clicks", Label: "CIMB Clicks"},
		{Value: "danamon_online", Label: "Danamon Online Banking"},
		{Value: "akulaku", Label: "Akulaku"},
		{Value: "kredivo", Label: "Kredivo"},
		{Value: "indomaret", Label: "Indomaret"},
		{Value: "alfamart", Label: "Alfamart"},
	},
	"xendit": {
		{Value: xendit.PaymentMethodCreditCard, Label: "Credit Card"},
		{Value: xendit.PaymentMethodBCA, Label: "BCA Virtual Account"},
		{Value: xendit.PaymentMethodBNI, Label: "BNI Virtual Account"},
		{Value: xendit.PaymentMethodBSI, Label: "BSI Virtual Account"},
		{Value: xendit.PaymentMethodBRI, Label: "BRI Virtual Account"},
		{Value: xendit.PaymentMethodMandiri, Label: "Mandiri Virtual Account"},
		{Value: xendit.PaymentMethodPermata, Label: "Permata Virtual Account"},
		{Value: xendit.PaymentMethodSahabatSampoerna, Label: "Sahabat Sampoerna"},
		{Value: xendit.PaymentMethodBNC, Label: "BNC Virtual Account"},
		{Value: xendit.PaymentMethodAlfamart, Label: "Alfamart"},
		{Value: xendit.PaymentMethodIndomaret, Label: "Indomaret"},
		{Value: xendit.PaymentMethodOVO, Label: "OVO"},
		{Value: xendit.PaymentMethodDana, Label: "DANA"},
		{Value: xendit.PaymentMethodShopeePay, Label: "ShopeePay"},
		{Value: xendit.PaymentMethodLinkAja, Label: "LinkAja"},
		{Value: xendit.PaymentMethodJeniusPay, Label: "Jenius Pay"},
		{Value: xendit.PaymentMethodDD_BRI, Label: "Direct Debit BRI"},
		{Value: xendit.PaymentMethodDD_BCA_KlikPay, Label: "Direct Debit BCA KlikPay"},
		{Value: xendit.PaymentMethodKredivo, Label: "Kredivo"},
		{Value: xendit.PaymentMethodAkulaku, Label: "Akulaku"},
		{Value: xendit.PaymentMethodAtome, Label: "Atome"},
		{Value: xendit.PaymentMethodQRIS, Label: "QRIS"},
	},
	"tripay": {
		{Value: tripay.ChannelPermataVA, Label: "Permata Virtual Account"},
		{Value: tripay.ChannelBniVA, Label: "BNI Virtual Account"},
		{Value: tripay.ChannelBriVA, Label: "BRI Virtual Account"},
		{Value: tripay.ChannelMandiriVA, Label: "Mandiri Virtual Account"},
		{Value: tripay.ChannelBcaVA, Label: "BCA Virtual Account"},
		{Value: tripay.ChannelMuamalatVA, Label: "Muamalat Virtual Account"},
		{Value: tripay.ChannelCimbVA, Label: "CIMB Niaga Virtual Account"},
		{Value: tripay.ChannelBsiVA, Label: "BSI Virtual Account"},
		{Value: tripay.ChannelOcbcVA, Label: "OCBC NISP Virtual Account"},
		{Value: tripay.ChannelDanamonVA, Label: "Danamon Virtual Account"},
		{Value: tripay.ChannelOtherBankVA, Label: "Other Bank Virtual Account"},
		{Value: tripay.ChannelAlfamart, Label: "Alfamart"},
		{Value: tripay.ChannelIndomaret, Label: "Indomaret"},
		{Value: tripay.ChannelAlfamidi, Label: "Alfamidi"},
		{Value: tripay.ChannelOvo, Label: "OVO"},
		{Value: tripay.ChannelQris, Label: "QRIS (ShopeePay)"},
		{Value: tripay.ChannelQrisc, Label: "QRIS (Customizable)"},
		{Value: tripay.ChannelQris2, Label: "QRIS 2"},
		{Value: tripay.ChannelDana, Label: "DANA"},
		{Value: tripay.ChannelShopeePay, Label: "ShopeePay"},
		{Value: tripay.ChannelQrisShopeePay, Label: "QRIS Custom (ShopeePay)"},
	},
	"duitku": {
		{Value: "va", Label: "Virtual Account"},
		{Value: "credit_card", Label: "Credit Card"},
		{Value: "qris", Label: "QRIS"},
		{Value: "shopeepay", Label: "ShopeePay"},
		{Value: "ovo", Label: "OVO"},
		{Value: "dana", Label: "DANA"},
		{Value: "linkaja", Label: "LinkAja"},
		{Value: "indomaret", Label: "Indomaret"},
		{Value: "alfamart", Label: "Alfamart"},
		{Value: "cimb_clicks", Label: "CIMB Clicks"},
	},
	"stripe": {
		{Value: "card", Label: "Card"},
		{Value: "bank_transfer", Label: "Bank Transfer"},
		{Value: "alipay", Label: "Alipay"},
		{Value: "wechat_pay", Label: "WeChat Pay"},
		{Value: "klarna", Label: "Klarna"},
		{Value: "afterpay", Label: "Afterpay"},
		{Value: "giropay", Label: "Giropay"},
		{Value: "ideal", Label: "iDEAL"},
		{Value: "sepa_debit", Label: "SEPA Direct Debit"},
		{Value: "sofort", Label: "SOFORT"},
	},
	"paypal": {
		{Value: "paypal", Label: "PayPal"},
		{Value: "credit_card", Label: "Credit Card"},
		{Value: "debit_card", Label: "Debit Card"},
		{Value: "venmo", Label: "Venmo"},
	},
	"other": {},
}

// providerMethodLabel returns the display label for a provider+method pair.
func providerMethodLabel(provider, value string) string {
	for _, m := range providerMethods[provider] {
		if m.Value == value {
			return m.Label
		}
	}
	return value
}

// resolveMethod parses "provider|value" into its parts and validates it.
func resolveMethod(s string) (provider string, value string, ok bool) {
	for p := range providerMethods {
		for _, m := range providerMethods[p] {
			if s == p+"|"+m.Value {
				return p, m.Value, true
			}
		}
	}
	return "", "", false
}

// fieldsFor returns the dynamic credential field definitions for a provider.
func fieldsFor(provider string) []MethodField {
	if f, ok := providerMethodFields[provider]; ok {
		return f
	}
	return []MethodField{}
}

func handleListMethods(c fiber.Ctx) error {
	out := make([]CatalogMethod, 0)
	for provider, methods := range providerMethods {
		fields := fieldsFor(provider)
		for _, m := range methods {
			out = append(out, CatalogMethod{
				Provider: provider,
				Value:    m.Value,
				Label:    m.Label,
				Fields:   fields,
			})
		}
	}
	return c.JSON(fiber.Map{"methods": out})
}
