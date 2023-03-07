Proof of concept for streaming message via http server-sent event as Websocket alternative for handling Half-Duplex Communication.

## Setup and Testing

1. Run App

```
go run .
```

2. Subscribe notification of an user

```curl
curl --location 'http://localhost/notifications' \
--header 'X-User-ID: YOUR_USER_ID'
```

3. Sending an notification to specific user

```curl
curl --location 'http:/localhost/notifications' \
--header 'Content-Type: application/json' \
--data '{
    "user_id": "YOUR_USER_ID",
    "message": "Hello There!"
}'
```

