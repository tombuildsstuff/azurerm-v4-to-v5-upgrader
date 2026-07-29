package provider

import "github.com/tombuildsstuff/azurerm-v4-to-v5-upgrader/internal/rules"

func init() {
	rules.Register(NewVersionPinRule())
	rules.Register(NewProviderRegistrationRule())
	rules.Register(NewEnhancedValidationRule())
}
