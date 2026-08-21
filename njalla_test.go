package njalla

import (
	"testing"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	libdnsnjalla "github.com/libdns/njalla"
)

func TestUnmarshalCaddyfile(t *testing.T) {
	cases := []struct {
		name      string
		config    string
		wantToken string
		wantErr   bool
	}{
		{
			name:      "inline argument",
			config:    "njalla abc123",
			wantToken: "abc123",
		},
		{
			name: "token in block",
			config: `njalla {
				api_token abc123
			}`,
			wantToken: "abc123",
		},
		{
			name:      "placeholder is preserved for provisioning",
			config:    "njalla {env.NJALLA_API_TOKEN}",
			wantToken: "{env.NJALLA_API_TOKEN}",
		},
		{
			name:    "no token at all",
			config:  "njalla",
			wantErr: true,
		},
		{
			name: "empty block",
			config: `njalla {
			}`,
			wantErr: true,
		},
		{
			name:    "too many arguments",
			config:  "njalla abc123 def456",
			wantErr: true,
		},
		{
			// Accepting both would mean discarding one of the two tokens the
			// user wrote, so it is rejected rather than silently ignored.
			name: "inline argument and a block",
			config: `njalla abc123 {
				api_token def456
			}`,
			wantErr: true,
		},
		{
			name: "unrecognized subdirective",
			config: `njalla {
				api_key abc123
			}`,
			wantErr: true,
		},
		{
			name: "subdirective without a value",
			config: `njalla {
				api_token
			}`,
			wantErr: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			p := Provider{new(libdnsnjalla.Provider)}
			err := p.UnmarshalCaddyfile(caddyfile.NewTestDispenser(tc.config))

			if tc.wantErr {
				if err == nil {
					t.Fatalf("expected an error, got token %q", p.Provider.APIToken)
				}
				return
			}
			if err != nil {
				t.Fatalf("UnmarshalCaddyfile: %v", err)
			}
			if p.Provider.APIToken != tc.wantToken {
				t.Errorf("APIToken = %q, want %q", p.Provider.APIToken, tc.wantToken)
			}
		})
	}
}

func TestProvisionResolvesPlaceholders(t *testing.T) {
	t.Setenv("NJALLA_TEST_TOKEN", "resolved-token")

	p := Provider{&libdnsnjalla.Provider{APIToken: "{env.NJALLA_TEST_TOKEN}"}}
	if err := p.Provision(caddy.Context{}); err != nil {
		t.Fatalf("Provision: %v", err)
	}
	if p.Provider.APIToken != "resolved-token" {
		t.Errorf("APIToken = %q, want the resolved value", p.Provider.APIToken)
	}
}

// An unset environment variable resolves to an empty token, which would
// otherwise fail later with a confusing authentication error.
func TestProvisionRejectsUnresolvedPlaceholder(t *testing.T) {
	p := Provider{&libdnsnjalla.Provider{APIToken: "{env.NJALLA_DEFINITELY_NOT_SET}"}}
	err := p.Provision(caddy.Context{})
	if err == nil {
		t.Fatal("expected an error for a placeholder that resolves to nothing")
	}
}

func TestCaddyModuleID(t *testing.T) {
	info := Provider{}.CaddyModule()
	if got := string(info.ID); got != "dns.providers.njalla" {
		t.Errorf("module ID = %q", got)
	}
	module, ok := info.New().(*Provider)
	if !ok {
		t.Fatalf("New() returned %T, want *Provider", info.New())
	}
	if module.Provider == nil {
		t.Error("New() must return a Provider with an embedded libdns provider")
	}
}
