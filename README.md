# distributed-key-value-store

A tiny distributed key-value store in Go: three participant servers hold email → password pairs **in memory**, and a coordinator routes each request to the right server using **consistent hashing**.

Data is lost when a server stops. This is a learning project, not production code.

## Architecture

```
                 ┌────────────────┐
  client ──────► │  coordinator   │  consistent hash ring, 307 redirect
  :3000          │  (main.go)     │──────────────┬──────────────┬─────────
                 └────────────────┘              ▼              ▼         ▼
                                            participant    participant  participant
                                             :8080          :8081        :8082
                                             (in-memory map)
```

- **coordinator** — listens on `:3000`, picks one of the 3 participant servers via a consistent hash ring and replies with an HTTP `307` redirect to it. The full path + query string are preserved.
- **participant** — listens on the port given as a CLI arg. Stores data in a plain Go map, so a key always lands on the same server.

Note: the hash is currently computed from the request's `Host` header, so all requests sent to `localhost:3000` land on the same node. Swapping in the email as the routing key is the next step.

## Run

Terminal 1 — coordinator:

```bash
cd coordinator
go run .
```

Terminals 2–4 — participants (one per port):

```bash
cd participants
go run main.go :8080
go run main.go :8081
go run main.go :8082
```

## Usage (through the coordinator on :3000)

```bash
# create a key (POST, JSON body)
curl -i -L http://localhost:3000/createUserKey \
  -H "Content-Type: application/json" \
  -d '{"Email":"souravpoddar6677@gmail.com","password":93042}'
# Key is Stored Sucessfully   (again: Key is already present)

# read a key (GET, email as a query param)
curl -i -L -G http://localhost:3000/getUserKey \
  --data-urlencode "EMAIL=souravpoddar6677@gmail.com"
# User Password is :93042   (unknown key: Key is not present first create)

# update a key (PUT, JSON body; key must exist first)
curl -i -L -X PUT http://localhost:3000/updateUserKey \
  -H "Content-Type: application/json" \
  -d '{"Email":"souravpoddar6677@gmail.com","password":999999}'
# User keys updated
```

`-L` follows the coordinator's redirect to the participant.

You can also hit a participant directly (`http://localhost:8080/...`) — it works, but you're bypassing the hash ring.

## Other routes

```bash
curl http://localhost:3000/home?name=Alice   # Hello, Alice!
curl http://localhost:3000/checkhealth       # memory stats of a participant
```
