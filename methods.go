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

func handleListMethods(c fiber.Ctx) error {
	// return a map: provider type -> list of methods (for the FE dropdown)
	out := make(fiber.Map, len(providerMethods))
	for provider, methods := range providerMethods {
		out[provider] = methods
	}
	return c.JSON(fiber.Map{"providers": out})
}

// validProviderMethod reports whether m is a known method for provider type t.
func validProviderMethod(t, m string) bool {
	for _, pm := range providerMethods[t] {
		if pm.Value == m {
			return true
		}
	}
	return false
}
