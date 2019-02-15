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

Meaning:

| Field   |      Type      | Must |  Description |
|----------|:-------------:|-----|------|
| user_name |  string | Y |your register unique user name |
| pass_wd |    string   | Y |  your password you set |
| remember | bool | N |   true will save cookie and 7 days login without login |

## Response

Normal:

```
{
  "flag": true,
  "cid": "fc067ace86c54f99b6084f044bcd31f5"
}
```

Meaning:

| Field   |      Type      |  Description |
|----------|:-------------:|------|
| flag |  bool | API request success will return true |
| cid |    string   |   unique id to debug API log |

Something wrong:

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

Meaning:

| Field   |      Type      |  Description |
|----------|:-------------:|------|
| flag |  bool | API request wrong will return false |
| cid |    string   |   unique id to debug API log |
| error | bool |    flag is false will return error: id and msg point the meaning |

All API will be return `flag` and `cid`, if `flag` is false, there will a error msg out.

## /logout

Logout will be clean auth cookie in client and clean session in server. 

```
POST /logout
```

will be always return:

```
{
  "flag": true
}
```

## /b/upload

Base API will prefix `/b`.

### Request

This API can request by HTML form `multipart/form-data`:

| Field   |      Type      | Must |  Description |
----------|:-------------:|-----|------|
| type |  string | N |limit type of file, can see below. "media" mean can only upload jpg, png... optional |
| describe |    string  | N |   upload describe can be empty |
| file | bin | Y |  file is file, max size 33.54MB |

type can be:

```
{
	"image": {
		"jpg", "jpeg", "png", "bmp", "gif"},
	"flash": {
		"swf", "flv"},
	"media": {
		"swf", "flv", "mp3", "wav", "wma", "wmv", "mid", "avi", "mpg", "asf", "rm", "rmvb"},
	"file": {
		"doc", "docx", "xls", "xlsx", "ppt", "htm", "html", "txt", "zip", "rar", "gz", "bz2", "pdf"},
	"other": {
		"jpg", "jpeg", "png", "bmp", "gif", "swf", "flv", "mp3",
		"wav", "wma", "wmv", "mid", "avi", "mpg", "asf", "rm", "rmvb",
		"doc", "docx", "xls", "xlsx", "ppt", "htm", "html", "txt", "zip", "rar", "gz", "bz2"}
}
```

### Response

All result will include in `data`:

Normal:

```
{
  "flag": true,
  "cid": "6a2bd7b2e1eb48efaafaf36fb07c8085",
  "data": {
    "file_name": "asr.wav",
    "size": 98124,
    "url": "/storage/media/-1/68756e746572687567510102e985d0f2f311ad3295534dc435.wav",
    "addon": "file the same in server"
  }
}
```

Meaning:

| Field   |      Type      |  Description |
|----------|:-------------:|------|
| file_name |  string | upload filename |
| size |  int | byte size of file |
| url |  string | inner url of this file, every file will has a unique md5, here is `68756e746572687567510102e985d0f2f311ad3295534dc435` |
| addon |  string | if file md5 the same will appear this |

Wrong:

```
{
  "flag": false,
  "cid": "f4e280f37de145a7a23a33ec7c4ccd11",
  "error": {
    "id": 10002,
    "msg": "upload file err http: no such file"
  }
}
```

## /v1/group/create

Admin API will has prefix `/v1`, and will be check auth every request.

### Request

```
{
	"name": "common user2w",
	"describe": "test group",
	"image_path": "/storage/media/-1/68756e746572687567510102e985d0f2f311ad3295534dc435.wav"
}
```

Meaning:

| Field   |      Type      | Must |  Description |
|----------|:-------------:|-----|------|
| name |  string | Y |group name, can not repeat |
| describe |    string   | N|  group describe, optional |
| image_path | bool | N |   image url of group, optional, if not empty must be the url return by /upload |

### Response

Normal:

```
{
  "flag": true,
  "cid": "de8c0b6b82694833872bf2065bc6e57d",
  "data": {
    "id": 13,
    "name": "common user2w",
    "describe": "test group",
    "create_time": 1550219797,
    "image_path": "/storage/media/-1/68756e746572687567510102e985d0f2f311ad3295534dc435.wav"
  }
}
```

Meaning:

| Field   |      Type      |  Description |
|----------|:-------------:|------|
| id |  string | group Id |
| create_time |  int | Unix timestamp of group create |

Wrong:

```
{
  "flag": false,
  "cid": "586c13cd8d1b4735ba1f1b406d0473ff",
  "error": {
    "id": 10005,
    "msg": "paras not right:group name exist"
  }
}

{
  "flag": false,
  "cid": "c26b18e8a9d3492090a14dcc9894e181",
  "error": {
    "id": 10005,
    "msg": "paras not right:name can not empty"
  }
}

{
  "flag": false,
  "cid": "2523072261974fd8a07a655e68b94ab5",
  "error": {
    "id": 10005,
    "msg": "paras not right:image not exist"
  }
}
```