# irccip

irccip implements Sony's InfraRed Compatible Control over Internet Protocol (IRCC-IP) for Bravia displays

## Usage

```go
// Create a client to send commands to a remote display
c := irccip.NewClient("http://192.168.1.12", "0000")

// Send the 'Home' key code
if err := c.SendKeyCode(irccip.KeyHome); err != nil {
    log.Fatalf("SendKeyCode error: %v", err)
}
```

## Key Codes

The package contains key code constants taken from a `KDL-43W809C` TV

### Finding new IRCC-IP key codes

Query your TV/Display, using the following command, to get a list of available key codes:

```
$ curl http://192.168.1.12/sony/system -d '{"method": "getRemoteControllerInfo", "id": 1, "params": [], "version": "1.0"}'
```
```
{
  "result": [
    {
      "bundled": true,
      "type": "IR_REMOTE_BUNDLE_TYPE_AEP_N"
    },
    [
      {
        "name": "Num1",
        "value": "AAAAAQAAAAEAAAAAAw=="
      },
      // ...
      {
        "name": "AndroidMenu",
        "value": "AAAAAgAAAMQAAABPAw=="
      }
    ]
  ],
  "id": 1
}
```
