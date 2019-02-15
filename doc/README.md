# API overview

you can use insomnia to debug API. All API will be request JSON and response JSON result.

## /login

Some API must auth, So login first.

### Request

```
POST /login

{
	"user_name": "hunterhug",
	"pass_wd": "fafa",
	"remember": true
}
```

| Field   |      Type      |  Note |
|----------|:-------------:|------:|
| user_name |  string | your register unique user name |
| pass_wd |    string   |   your password you set |
| remember | bool |    true will save cookie and 7 days login without login |

## Response

```
{
  "flag": false,
  "cid": "01a786b82ef847379d2a2a93b52146eb",
  "error": {
    "id": 10003,
    "msg": "username or password wrong"
  }
}
```