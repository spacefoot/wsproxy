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

Report the weight to be unstable
```json
{
  "type": "unstable",
  "data": {}
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

Request the last weight
```json
{
  "type": "weight"
}
```

Configure the logger
```json
{
  "type": "log",
  "data": {
    "enabled": true
  }
}
```

Debug only, send a custom weight to be broadcasted
```json
{
  "type": "debug-weight",
  "data": {
    "weight": 0,
    "unit": "g"
  }
}
```

Debug only, send unstable to be broadcasted
```json
{
  "type": "debug-unstable"
}
```
