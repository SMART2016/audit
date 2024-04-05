## Simple user Service
### Supported API's
- GET localhost:8080/user-service/v1/health
- GET localhost:8080/user-service/v1/users/{user_id}
- DELETE localhost:8080/user-service/v1/users/{user_id}
- PUT localhost:8080/user-service/v1/users/{user_id}
- POST localhost:8080/user-service/v1/users
  - The user payload is as below:
    `{
        "name":"",
        "email":"",
        "id":""
      }`