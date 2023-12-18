# Cart-Servie with go!

## Overview
- Midterm project of the course "Web Programming" at AUT
- A project that simulates an online shop cart handler written in Go and with RESTful API design.
- JWT is used for authentication and data is kept in a PostgreSQL database, and Gorm is used as the ORM.
- The data stored in database is in both normal format and also as jsonb data.

## Commands
You can use any tools like `curl` or `postman` to send api requests. Notice that you need to signup and login afterwards in order to use the subsequent commands like `GetAllBaskets` or `CreateBasket`.

**Signup**: First you need to signup.

- Using curl: `curl -X POST -H "Content-Type: application/json" -d '{"username": "your_username", "password": "your_password"}' localhost:2023/signup`. This will signup a new user if the username already doesn't exist.

**Login**: After signing up, you need to login again to get a JWT token, so you can continue with sending other requests. Notice that the given token will be valid for 5 minutes by default, but can be configured using `expDur` variable in `main.go` file.

- Using curl: `curl -X POST -H "Content-Type: application/json" -d '{"username": "your_username", "password": "your_password"}' localhost:2023/login`

**CreateBasket**: After a successful login, you will receive a message in which you can find your JWT token with the tag `token`. You have to send that token in any commands except `Signup` and `Login`.

- Using curl: `curl -X POST -H "Content-Type: application/json" -d '{"data": "data_to_be_stored_in_this_basket", "token": "given_token_from_login_step"}' localhost:2023/basket/`. This will create a basket with the given data for the owner of the token. Notice that for simplicity, each basket holds a `data` that is a string of at most 2048 bytes as a representation of the products stored in the basket. Also each basket has a field named `state` that represents the current state of the basket with can be `COMPLETED` or `PENDING`. When creating a new basket, server will set the state of the basket as `COMPLETED` by default.

**GetAllBaskets**: This command retrieves all the baskets created for the user, or if no basket exists returns a `Not Found` error.

- Using curl: `curl -H "Content-Type: application/json" -d '{"token": "given_token_from_login_step"}' localhost:2023/basket/`.

**GetBasketByID**: This command retrieves a specific basket of the user with the given `id`.

- Using curl: `curl -H "Content-Type: application/json" -d '{"token": "given_token_from_login_step"}' localhost:2023/basket/id_of_basket`

**UpdateBasket**: This command is uesd for updating two fields of the basket: `data` and `state`. Notice that giving both of these fields is arbitrary. Also notice that the state can be changed only if it hadn't been changed to `COMPLETED` previously.

- Using curl: `curl -X PATCH -H "Content-Type: application/json" -d '{"data": "new_data_of_basket", "state": "COMPLETED", "token": "given_token_from_login_step"}' localhost:2023/basket/id_of_basket`

**DeleteBasket**: This command is used for deleting a basket by its id.

- Using curl: `curl -X DELETE -H "Content-Type: application/json" -d '{"token": "given_token_from_login_step"}' localhost:2023/basket/id_of_basket`