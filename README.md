# Njalla module for Caddy

This package contains a DNS provider module for [Caddy](https://github.com/caddyserver/caddy).
It can be used to manage DNS records with [Njalla](https://njal.la) accounts,
and to solve the ACME DNS-01 challenge.

## Caddy module name

```
dns.providers.njalla
```

## Building

This module is not part of the standard Caddy distribution. Build Caddy with it
using [xcaddy](https://github.com/caddyserver/xcaddy):

```sh
xcaddy build --with github.com/caddy-dns/njalla
```

A `Dockerfile` is included for building a Caddy image with the module baked in:

```sh
docker build -t caddy-njalla .
```

## Authenticating

Create an API token from your Njalla account settings, under
[Settings → API](https://njal.la/settings/api/).

For the ACME DNS-01 challenge, Njalla's token editor offers an **ACME** option
that grants exactly what is needed: `add-record`, `edit-record`,
`remove-record`, and `list-records`, restricted to TXT records under the
`_acme-challenge` prefix.

See the [libdns provider README](https://github.com/libdns/njalla) for the full
list of permissions and for provider-specific behaviour.

## Config examples

### Caddyfile

Using the [`tls` directive](https://caddyserver.com/docs/caddyfile/directives/tls)
with the DNS challenge:

```caddy
tls {
	dns njalla {env.NJALLA_API_TOKEN}
}
```

The token may also be given in a block:

```caddy
tls {
	dns njalla {
		api_token {env.NJALLA_API_TOKEN}
	}
}
```

To apply it to every site, set it in the
[global options](https://caddyserver.com/docs/caddyfile/options):

```caddy
{
	acme_dns njalla {env.NJALLA_API_TOKEN}
}
```

Putting the token in an environment variable keeps it out of your Caddyfile.
Any [Caddy placeholder](https://caddyserver.com/docs/conventions#placeholders)
is resolved when the module is provisioned; a literal token works too.

### JSON

To use this module for the ACME DNS challenge,
[configure the ACME issuer](https://caddyserver.com/docs/json/apps/tls/automation/policies/issuer/acme/)
like so:

```json
{
  "module": "acme",
  "challenges": {
    "dns": {
      "provider": {
        "name": "njalla",
        "api_token": "YOUR_NJALLA_API_TOKEN"
      }
    }
  }
}
```

## Troubleshooting

**`API token is empty after resolving placeholders`** — the placeholder you used
resolved to nothing, usually because the environment variable is not set in the
process running Caddy. Note that a systemd service does not inherit your shell's
environment; set it with `Environment=` or an `EnvironmentFile=` in the unit.

**Certificates fail with a DNS challenge error** — confirm the token is allowed
to call `add-record` and `remove-record` for the `_acme-challenge` prefix. The
ACME option in Njalla's token editor sets this up for you.

**Wildcard certificates** — these require the DNS challenge, which is what this
module provides. Make sure the zone you are issuing for is the Njalla domain
itself, not a delegated subdomain.
