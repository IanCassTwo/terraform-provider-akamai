package providers

import (
	// This is where providers are import so they can register themselves
	_ "github.com/akamai/terraform-provider-akamai/v2/pkg/providers/dns"
	_ "github.com/akamai/terraform-provider-akamai/v2/pkg/providers/gtm"
	_ "github.com/akamai/terraform-provider-akamai/v2/pkg/providers/property"
)
