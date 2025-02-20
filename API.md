# API

All communication is done using websocket at `ws://localhost:23193/ws`.

All messages are JSON encoded, as a `type` field and an optional `data` field.

## From Server

Report the stable weight on the scale
```json
{
  "type": "weight",
  "data": {
    "weight": 0,
    "unit": "g"
  }
}
```

Report the status of the scale
```json
{
  "type": "status",
  "data": {
    "open": true
  }
}
```

## To Server

Request the current status
```json
{
  "type": "status"
}
```

Reset the scale value to 0
```json
{
  "type": "zero"
}
```

Debug only, send a custom weight to be broadcasted
```json
{
  "type": "weight",
  "data": {
    "weight": 0,
    "unit": "g"
  }
}
```
