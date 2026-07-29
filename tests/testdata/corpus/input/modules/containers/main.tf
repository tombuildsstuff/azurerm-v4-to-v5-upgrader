# ---------------------------------------------------------------------------
# Container Apps: environment + app + job
# ---------------------------------------------------------------------------

resource "azurerm_container_app_environment" "example" {
  name                = "example-cae"
  location            = "westeurope"
  resource_group_name = "example"
}

resource "azurerm_container_app" "example" {
  name                         = "example-ca"
  container_app_environment_id = azurerm_container_app_environment.example.id
  resource_group_name          = "example"
  revision_mode                = "Single"

  template {
    container {
      name   = "example"
      image  = "nginx:latest"
      cpu    = 0.25
      memory = "0.5Gi"

      liveness_probe {
        transport                        = "HTTP"
        port                             = 80
        termination_grace_period_seconds = 10
      }

      startup_probe {
        transport                        = "HTTP"
        port                             = 80
        termination_grace_period_seconds = 10
      }
    }
  }
}

resource "azurerm_container_app_job" "example" {
  name                         = "example-caj"
  location                     = "westeurope"
  container_app_environment_id = azurerm_container_app_environment.example.id
  resource_group_name          = "example"
  replica_timeout_in_seconds   = 10

  manual_trigger_config {
    parallelism              = 1
    replica_completion_count = 1
  }

  template {
    container {
      name   = "example"
      image  = "nginx:latest"
      cpu    = 0.25
      memory = "0.5Gi"

      liveness_probe {
        transport                        = "HTTP"
        port                             = 80
        termination_grace_period_seconds = 10
      }

      startup_probe {
        transport                        = "HTTP"
        port                             = 80
        termination_grace_period_seconds = 10
      }
    }
  }
}

# ---------------------------------------------------------------------------
# Container Registry
# ---------------------------------------------------------------------------

resource "azurerm_container_registry" "example" {
  name                = "examplecr"
  resource_group_name = "example"
  location            = "westeurope"
  sku                 = "Premium"

  trust_policy_enabled = true

  georeplications {
    location                  = "northeurope"
    regional_endpoint_enabled = true
  }
}

# ---------------------------------------------------------------------------
# Classic CDN: profile + endpoint + custom domain
# ---------------------------------------------------------------------------

resource "azurerm_cdn_profile" "example" {
  name                = "example-cdn-profile"
  location            = "westeurope"
  resource_group_name = "example"
  sku                 = "Standard_Microsoft"
}

resource "azurerm_cdn_endpoint" "example" {
  name                = "example-cdn-endpoint"
  profile_name        = azurerm_cdn_profile.example.name
  location            = "westeurope"
  resource_group_name = "example"

  origin {
    name      = "example-origin"
    host_name = "www.contoso.com"
  }
}

resource "azurerm_cdn_endpoint_custom_domain" "example" {
  name            = "example"
  cdn_endpoint_id = azurerm_cdn_endpoint.example.id
  host_name       = "example.contoso.com"

  cdn_managed_https {
    certificate_type = "Dedicated"
    protocol_type    = "ServerNameIndication"
    tls_version      = "TLS12"
  }
}

resource "azurerm_cdn_endpoint_custom_domain" "example2" {
  name            = "example2"
  cdn_endpoint_id = azurerm_cdn_endpoint.example.id
  host_name       = "example2.contoso.com"

  user_managed_https {
    key_vault_secret_id = "https://example.vault.azure.net/secrets/cert/abc123"
    tls_version         = "TLS12"
  }
}

# ---------------------------------------------------------------------------
# Front Door (Standard/Premium): profile + origin group + rule set + rule +
# custom domain + WAF policy + security policy
# ---------------------------------------------------------------------------

resource "azurerm_cdn_frontdoor_profile" "example" {
  name                = "example-afd-profile"
  resource_group_name = "example"
  sku_name            = "Standard_AzureFrontDoor"
}

resource "azurerm_cdn_frontdoor_origin_group" "example" {
  name                     = "example-og"
  cdn_frontdoor_profile_id = azurerm_cdn_frontdoor_profile.example.id

  load_balancing {
    sample_size                 = 4
    successful_samples_required = 3
  }
}

resource "azurerm_cdn_frontdoor_rule_set" "example" {
  name                     = "examplers"
  cdn_frontdoor_profile_id = azurerm_cdn_frontdoor_profile.example.id
}

resource "azurerm_cdn_frontdoor_custom_domain" "example" {
  name                     = "example"
  cdn_frontdoor_profile_id = azurerm_cdn_frontdoor_profile.example.id
  host_name                = "example.contoso.com"

  tls {
    certificate_type    = "ManagedCertificate"
    minimum_tls_version = "TLS12"
  }
}

resource "azurerm_cdn_frontdoor_rule" "example" {
  name                      = "examplerule"
  cdn_frontdoor_rule_set_id = azurerm_cdn_frontdoor_rule_set.example.id
  order                     = 1
  behavior_on_match         = "Continue"

  actions {
    request_header_action {
      header_action = "Append"
      header_name   = "X-A"
      value         = "1"
    }
    response_header_action {
      header_action = "Overwrite"
      header_name   = "X-B"
      value         = "2"
    }
    url_redirect_action {
      redirect_type        = "Found"
      destination_hostname = "contoso.com"
    }
    url_rewrite_action {
      source_pattern          = "/a"
      destination             = "/b"
      preserve_unmatched_path = false
    }
    route_configuration_override_action {
      cdn_frontdoor_origin_group_id = azurerm_cdn_frontdoor_origin_group.example.id
      forwarding_protocol           = "HttpsOnly"
      cache_behavior                = "OverrideAlways"
      cache_duration                = "1.12:00:00"
      compression_enabled           = true
      query_string_caching_behavior = "IgnoreQueryString"
    }
  }

  conditions {
    client_port_condition {
      operator     = "Equal"
      match_values = ["8080"]
    }
    cookies_condition {
      cookie_name  = "session"
      operator     = "Equal"
      match_values = ["abc"]
    }
    host_name_condition {
      operator         = "Equal"
      match_values     = ["contoso.com"]
      negate_condition = true
    }
    http_version_condition {
      operator     = "Equal"
      match_values = ["2.0"]
    }
    is_device_condition {
      operator     = "Equal"
      match_values = ["Mobile"]
    }
    post_args_condition {
      post_args_name = "field"
      operator       = "Equal"
      match_values   = ["value"]
    }
    query_string_condition {
      operator     = "Contains"
      match_values = ["foo=bar"]
    }
    remote_address_condition {
      operator     = "IPMatch"
      match_values = ["10.0.0.0/8"]
    }
    request_body_condition {
      operator     = "Contains"
      match_values = ["payload"]
    }
    request_header_condition {
      header_name  = "X-Custom"
      operator     = "Equal"
      match_values = ["value"]
    }
    request_method_condition {
      operator     = "Equal"
      match_values = ["GET"]
    }
    request_scheme_condition {
      operator     = "Equal"
      match_values = ["HTTPS"]
    }
    request_uri_condition {
      operator     = "Contains"
      match_values = ["/api"]
    }
    server_port_condition {
      operator     = "Equal"
      match_values = ["443"]
    }
    socket_address_condition {
      operator     = "IPMatch"
      match_values = ["10.0.0.0/8"]
    }
    ssl_protocol_condition {
      operator     = "Equal"
      match_values = ["TLSv1.2"]
    }
    url_file_extension_condition {
      operator     = "Equal"
      match_values = ["html"]
    }
    url_filename_condition {
      operator     = "Equal"
      match_values = ["index"]
    }
    url_path_condition {
      operator     = "Contains"
      match_values = ["/blog"]
    }
  }
}

resource "azurerm_cdn_frontdoor_firewall_policy" "example" {
  name                = "examplewaf"
  resource_group_name = "example"
  sku_name            = azurerm_cdn_frontdoor_profile.example.sku_name
  mode                = "Prevention"
}

resource "azurerm_cdn_frontdoor_security_policy" "example" {
  name                     = "examplesp"
  cdn_frontdoor_profile_id = azurerm_cdn_frontdoor_profile.example.id

  security_policies {
    firewall {
      cdn_frontdoor_firewall_policy_id = azurerm_cdn_frontdoor_firewall_policy.example.id

      association {
        patterns_to_match = ["/*"]

        domain {
          cdn_frontdoor_domain_id = azurerm_cdn_frontdoor_custom_domain.example.id
        }
      }
    }
  }
}

# ---------------------------------------------------------------------------
# NGINX deployment: delegated subnet + public IP + deployment
# ---------------------------------------------------------------------------

resource "azurerm_virtual_network" "nginx" {
  name                = "example-nginx-vnet"
  address_space       = ["10.1.0.0/16"]
  location            = "westeurope"
  resource_group_name = "example"
}

resource "azurerm_subnet" "nginx" {
  name                 = "example-nginx-subnet"
  resource_group_name  = "example"
  virtual_network_name = azurerm_virtual_network.nginx.name
  address_prefixes     = ["10.1.1.0/24"]

  delegation {
    name = "delegation"

    service_delegation {
      name    = "NGINX.NGINXPLUS/nginxDeployments"
      actions = ["Microsoft.Network/virtualNetworks/subnets/join/action"]
    }
  }
}

resource "azurerm_public_ip" "nginx" {
  name                = "example-nginx-pip"
  resource_group_name = "example"
  location            = "westeurope"
  allocation_method   = "Static"
  sku                 = "Standard"
}

resource "azurerm_nginx_deployment" "example" {
  name                = "example-nginx"
  resource_group_name = "example"
  location            = "westeurope"
  sku                 = "standard_Monthly"

  diagnose_support_enabled = true
  managed_resource_group   = "example-managed-rg"

  frontend_public {
    ip_address = [azurerm_public_ip.nginx.id]
  }

  network_interface {
    subnet_id = azurerm_subnet.nginx.id
  }
}
