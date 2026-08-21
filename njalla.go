// Package njalla provides a Caddy DNS provider module for Njalla, wrapping
// the libdns implementation at github.com/libdns/njalla.
package njalla

import (
	"fmt"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	libdnsnjalla "github.com/libdns/njalla"
)

// Provider wraps the Njalla libdns provider as a Caddy module.
type Provider struct{ *libdnsnjalla.Provider }

func init() {
	caddy.RegisterModule(Provider{})
}

// CaddyModule returns the Caddy module information.
func (Provider) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "dns.providers.njalla",
		New: func() caddy.Module { return &Provider{new(libdnsnjalla.Provider)} },
	}
}

// Provision resolves placeholders in the API token before it is used.
// Implements caddy.Provisioner.
func (p *Provider) Provision(ctx caddy.Context) error {
	p.Provider.APIToken = caddy.NewReplacer().ReplaceAll(p.Provider.APIToken, "")
	if p.Provider.APIToken == "" {
		return fmt.Errorf("njalla: API token is empty after resolving placeholders; " +
			"check that the referenced environment variable is set")
	}
	return nil
}

// UnmarshalCaddyfile sets up the DNS provider from Caddyfile tokens. Two
// syntaxes are supported:
//
//	njalla <api_token>
//
//	njalla {
//	    api_token <api_token>
//	}
//
// Placeholders such as {env.NJALLA_API_TOKEN} are resolved during provisioning.
func (p *Provider) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {
	d.Next() // consume the directive name

	if d.NextArg() {
		p.Provider.APIToken = d.Val()
	} else {
		for nesting := d.Nesting(); d.NextBlock(nesting); {
			switch d.Val() {
			case "api_token":
				if !d.NextArg() {
					return d.ArgErr()
				}
				p.Provider.APIToken = d.Val()
			default:
				return d.Errf("unrecognized subdirective '%s'", d.Val())
			}
		}
	}

	if d.NextArg() {
		return d.Errf("unexpected argument '%s'", d.Val())
	}
	if p.Provider.APIToken == "" {
		return d.Err("missing API token")
	}
	return nil
}

// Interface guards
var (
	_ caddyfile.Unmarshaler = (*Provider)(nil)
	_ caddy.Provisioner     = (*Provider)(nil)
)
