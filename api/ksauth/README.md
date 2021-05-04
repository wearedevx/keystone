# Keystone Auth Server

It is used to auth users using their GitHub identity.

## Routes

### `POST /login-request`

> Starts a login procedure by creating a login request
> It will be used to retrieve user data when the login
> is successful on the third party (i.e. GitHub), and
> the user granted access to the application.

Request:  
No body is expected in the request.

Response:  
The `temporary_code` should be used as the `code` parameter in
the `GET /login-request` request.

```
Content-Type: application/json; charset=utf-8

{
   "id": 12,
   "temporary_code": "a temporary code as a string",
   "auth_code": "a auth code",
   "answered": true,
   "created_at": "2021-07-23T00:00:00.000Z",
   "updated_at": "2021-07-23T00:00:00.000Z"
}
```

### `GET /login-request`

> Should be polled to check the login status.

Request:  
Query string parameters:  
`code`: a temporay code to retrieve the login request

Response:

```
Content-Type: application/json; charset=utf-8

{
  "id": 12,
  "temporary_code": "a temporay code as a string",
  "auth_code": "a auth code",
  "answered": true,
  "created_at": "2021-07-23T00:00:00.000Z",
  "updated_at": "2021-07-23T00:00:00.000Z"
}
```

### `GET /auth-redirect/:code`

> Used by the third-party OAuth2 login process.
> When the user logs into their GitHub account and grants
> access to the app, the route is called, and the
> `LoginRequest` associated with `code` is updated.

### `POST /complete`

> Completes the login process.

Request:  
Body:

```
Content-Type: application/json; charset=utf-8

{
  "account_type": "github",
  "token": "eoiruer432dcb wdfkle w94d8",
  "public_key": "the user public key",
}
```

The `token` will be used to call the third’s party API to retreive
user informations such as their name, email and so on.

The `public_key` will be used to encrypt commication with the user.

Response:

```
Content-Type: application/json; charset=utf-8

{
  "id": 23,        // in databas id
  "ext_id": "wdf", // third party’s identifier
  "user_id": "aeruoiaur-eroieru-eoruise-seoriuesro", // unique identifier
  "account_type": "github",
  "username": "gaelph",
  "fullnamd": "Les Licornes ont des yeux",
  "email": "licornes@ont-des-yeux.com",
  "keys": {
    "sign": "public key to sign things",
    "ciper": "public ket to encrypt things"
  }
}
```
