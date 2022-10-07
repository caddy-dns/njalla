# Njalla module for Caddy

This package contains a DNS provider module for [Caddy](https://github.com/caddyserver/caddy). It can be used to manage DNS records for Njalla accounts.

## Caddy module name

```
dns.providers.njalla
```

## Config examples

To use this module for the ACME DNS challenge, [configure the ACME issuer in your Caddy JSON](https://caddyserver.com/docs/json/apps/tls/automation/policies/issuer/acme/) like so:

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

## Authenticating

See [the associated README in the libdns package](https://github.com/libdns/njalla) for important information about credentials.
