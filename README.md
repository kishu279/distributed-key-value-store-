# distributed-key-value-store

A small Go server that stores emails and passwords **in memory** (data is lost when it stops).

## Run

```bash
go run .
```

Server runs at `http://localhost:8080`.

## Routes

Key routes take a JSON body with `email` (string) and `password` (number):

```bash
# create
curl -X POST http://localhost:8080/createUserKey \
  -d '{"email":"alice@example.com","password":123456}'
# Key is Stored Sucessfully   (again: Key is already present)

# update
curl -X PUT http://localhost:8080/updateUserKey \
  -d '{"email":"alice@example.com","password":999999}'
# User keys updated   (unknown key: Key is not present first create)

# read (GET still needs a JSON body)
curl -X GET http://localhost:8080/getUserKey \
  -d '{"email":"alice@example.com","password":0}'
# User Password is :999999   (unknown key: Key is not present first create)
```

Other routes:

```bash
curl http://localhost:8080/              # Welcome to the Go Web Server!
curl "http://localhost:8080/home?name=Alice"   # Hello, Alice!
curl http://localhost:8080/checkhealth   # Total: ..., Available: ..., UsedPercent:...%
```
