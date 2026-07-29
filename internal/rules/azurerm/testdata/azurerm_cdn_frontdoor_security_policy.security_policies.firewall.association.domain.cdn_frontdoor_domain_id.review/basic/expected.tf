resource "azurerm_cdn_frontdoor_security_policy" "example" {
  name                     = "example"
  cdn_frontdoor_profile_id = "/subscriptions/x/profiles/p"

  security_policies {
    firewall {
      cdn_frontdoor_firewall_policy_id = "/subscriptions/x/policies/wafp"

      association {
        patterns_to_match = ["/*"]

        domain {
          cdn_frontdoor_domain_id = "/subscriptions/x/customDomains/cd"
        }
      }
    }
  }
}
