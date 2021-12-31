# caddy-jwt

![Go Workflow](https://github.com/ggicci/caddy-jwt/actions/workflows/go.yml/badge.svg) [![codecov](https://codecov.io/gh/ggicci/caddy-jwt/branch/main/graph/badge.svg?token=4V9OX8WFAW)](https://codecov.io/gh/ggicci/caddy-jwt) [![Go Reference](https://pkg.go.dev/badge/github.com/ggicci/caddy-jwt.svg)](https://pkg.go.dev/github.com/ggicci/caddy-jwt)

A Caddy HTTP Module - who Facilitates **JWT Authentication**

This module fulfilled [`http.handlers.authentication`](https://caddyserver.com/docs/modules/http.handlers.authentication) middleware as a provider named `jwt`.

[Documentation](https://caddyserver.com/docs/modules/http.authentication.providers.jwt)

## Install

Build this module with `caddy` at Caddy's official [download](https://caddyserver.com/download) site. Or:

```bash
xcaddy --with github.com/ggicci/caddy-jwt
```

## Quick View

```bash
git clone https://github.com/ggicci/caddy-jwt.git
cd caddy-jwt

# Build a caddy with this module and run an example server at localhost.
make example

TEST_TOKEN=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NTU4OTI2NzAsImp0aSI6IjgyMjk0YTYzLTk2NjAtNGM2Mi1hOGE4LTVhNjI2NWVmY2Q0ZSIsInVpZCI6MzQwNjMyNzk2MzUxNjkzMiwidXNlcm5hbWUiOiJnZ2ljY2kiLCJuc2lkIjozNDA2MzMwMTU3MTM3OTI2fQ.HWHw4qX4OGgCyNNa5En_siktjpoulTNwABXpEwQI4Q8

curl -v "http://localhost:8080?access_token=${TEST_TOKEN}"
# You should see authenticated output:
#
# User Authenticated with ID: 3406327963516932
#
# And the following command should also work:
curl -v -H"X-Api-Token: ${TEST_TOKEN}" "http://localhost:8080"
curl -v -H"Authorization: Bearer ${TEST_TOKEN}" "http://localhost:8080"
```

**NOTE**: you can decode the `${TEST_TOKEN}` above at [jwt.io](https://jwt.io/) to get human readable payload as follows:

```json
{
  "exp": 1655892670,
  "jti": "82294a63-9660-4c62-a8a8-5a6265efcd4e",
  "uid": 3406327963516932,
  "username": "ggicci",
  "nsid": 3406330157137926
}
```

## Configurations

Sample configuration (find more under [example](./example)):

```Caddyfile
api.example.com {
	route * {
		jwtauth {
			sign_key TkZMNSowQmMjOVU2RUB0bm1DJkU3U1VONkd3SGZMbVk=
			from_query access_token token
			from_header X-Api-Token
			from_cookies user_session
			issuer_whitelist https://api.example.com
			audience_whitelist https://api.example.io https://learn.example.com
			user_claims aud uid user_id username login
			meta_claims "IsAdmin->is_admin"
		}
		reverse_proxy http://172.16.0.14:8080
	}
}
```

**NOTE**:

1. Use `base64` to encode your key in the configuration.
2. The priority of `from_xxx` is `from_query > from_header > from_cookies`.

This module behaves like a "JWT Validator". Who

1. Extract the token from cookies, header or query from the HTTP request.
2. Validate the token by using the `sign_key`.
3. If the token is invalid by any reason, auth **failed** with `401`. Otherwise, next.
4. Get user id by inspecting the claims defined by `user_claims`.
5. If no valid user id (non-empty string) found, auth **failed** with `401`. Otherwise, next.
6. Return the user id to Caddy's authentication handler, and the context value `{http.auth.user.id}` got set. If `meta_claims` defined, user metadata placeholders `{http.auth.user.*}` will be populated, too.

## JWT Resources

- **MUST READ**: [JWT Security Best Practices](https://curity.io/resources/learn/jwt-best-practices/)
- Online Debuger: http://jwt.io/
