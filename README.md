# Mimoto

User authentication backend. Allows for registration and confirmation via email, and password resets via email.

Uses JWTs and refresh tokens.


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


- [ ] separate auth and user service concerns
- [ ] catch typed errors, create custom errors, and use appropriate messages
- [ ] cleanup common request/response formats
- [ ] rearchitect for unit testing
- [ ] setup logging