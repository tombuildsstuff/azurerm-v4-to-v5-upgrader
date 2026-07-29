// Package azurerm blank-imports every service rule package so that importing
// it populates the global rules registry.
package azurerm

import (
	_ "github.com/tombuildsstuff/azurerm-v4-to-v5-upgrader/internal/rules/azurerm/attr"
	_ "github.com/tombuildsstuff/azurerm-v4-to-v5-upgrader/internal/rules/azurerm/provider"
	_ "github.com/tombuildsstuff/azurerm-v4-to-v5-upgrader/internal/rules/azurerm/resource"
)
