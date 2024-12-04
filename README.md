# JWT Authentication and Authorization service

- JWT
- PostgreSQL (GORM)
- Golang + Gin

## Database
- users
  - id
  - email
####
- refresh_tokens
  - id
  - user_id
  - token_hashed
  - user_ip
  - expires
  - created_at

## Setup:

## Environment Variables Configuration

```
JWT_SECRET_KEY=your secret JWT auth key
DB_SERVER=localhost
DB_NAME=your database name
DB_PORT=5432
DB_USERNAME=postgres
DB_USER_PASSWORD=password
```

## DB:
```
CREATE DATABASE your_db_name;
```

### Add user for tests if needed
### main.go:
```
err := repository.AddUser("some-mail@email.com")
if err != nil {
	return 
}
```

## Access token lifespan - 15 minutes
## Refresh token lifespan - 48 hours

## Available routes:
## POST /generate_tokens
### Generates access and refresh tokens pair for the given user data: user_id, user_ip
- Sample request:
```
{
  "user_id": "c0496915-c052-4e89-b192-d67696a8b26b",
  "user_ip": "127.0.0.1"
}
```
- Sample response:
```
{
    "access_token": "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MzMzMzMxNTcsInVzZXJfaWQiOiJjMDQ5NjkxNS1jMDUyLTRlODktYjE5Mi1kNjc2OTZhOGIyNmIiLCJ1c2VyX2lwIjoiMTI3LjAuMC4xIn0.NxmzT2fXTNnZ23ekFj85wLpRTZ6lYw2IceS2QhB-45oG3qNukfgOEloyYOhtRE5Cl_mh-y0tX6OnIWB2uipQDA",
    "refresh_token": "Fl0UbQAsZuNhyA9b+xG18tahjmhcW4VcxdFF+/lzOSk="
}
```

## POST /refresh_tokens
### Refreshes the tokens for the given user data and unhashed refresh_token
- Sample request:
```
{
  "refresh_token": "Fl0UbQAsZuNhyA9b+xG18tahjmhcW4VcxdFF+/lzOSk=",
  "user_ip": "127.0.0.1",
  "user_id": "c0496915-c052-4e89-b192-d67696a8b26b"
}
```
- Sample response:
```
{
    "access_token": "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MzMzMzMxNzAsInVzZXJfaWQiOiJjMDQ5NjkxNS1jMDUyLTRlODktYjE5Mi1kNjc2OTZhOGIyNmIiLCJ1c2VyX2lwIjoiMTI3LjAuMC4yIn0.VLSdsDPfS1J_urK7bw0Mh4dEhbfqlc1veN1zYxEFRsdcDRRNH7c1DoiO8r0fpPg3Zt8I81MNJxvWhiRH2nWIgg",
    "refresh_token": "rpvq8wjoKQTRDW1PBRHNgaUdIRxmtkLrIV8xfGW9k7Y="
}
```
### Sends mock email to user if IP address has changed
