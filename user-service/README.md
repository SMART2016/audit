## Simple user Service
### Supported API's
- GET localhost:8080/user-service/v1/health
  - curl http://localhost:8080/user-service/v1/health
- GET localhost:8080/user-service/v1/users/{user_id}
  - curl http://localhost:8080/user-service/v1/users/1
- DELETE localhost:8080/user-service/v1/users/{user_id}
  - curl -X DELETE http://localhost:8080/user-service/v1/users/1
- PUT localhost:8080/user-service/v1/users/{user_id}
  - curl -X PUT -H "Content-Type: application/json" -d '{"id":"1", "name":"<newName>", "email":"<newEmail>"}' http://localhost:8080/user-service/v1/users/1
- POST localhost:8080/user-service/v1/users
  - curl -X POST -H "Content-Type: application/json" -d '{"id":"1", "name":"<name>", "email":"<email>"}' http://localhost:8080/user-service/v1/users
