# Mimoto

User authentication backend. Allows for registration and confirmation via email, and password resets via email.

Uses JWTs and refresh tokens.

`POST /signup`
```json
{
  "email": "test@test.com",
  "name": "test",
  "password": "secret"
}
```

`POST /confirm?email=test@test.com&code=jd73hjd-sd73h3-kj56d-sdf898`

`POST /login`
```json
{
  "email": "test@test.com",
  "password": "secret"
}
```

`POST /refresh`

`POST /sendreset`
```json
{
  "email": "test@test.com"
}
```

`POST /reset?email=test@test.com&code=2h48fj-f93jk3d-sdf987-lk6j7`
```json
{
  "password": "secret"
}
```

`POST /logout`





### TODO
- [x] structure project
- [x] setup chi-router server and routes
- [x] setup gorm
- [x] setup docker-compose for postgres
- [x] implement signup endpoint
- [x] implement confirm account endpoint
- [x] implement login endpoint
- [x] implement refresh endpoint
- [x] implement logout endpoint
- [x] implement reset endpoint

- [ ] catch typed errors, create custom errors, and use approeriate messages
- [x] cleanup common request/response formats
- [x] rearchitect for unit testing
- [ ] setup logging