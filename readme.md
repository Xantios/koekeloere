# Filewatcher
Piertje en moffel koekeloere naar je files

Basic filewatcher that calls a *webhook* when a change is noticed

# Build
Just run the included makefile without any arguments

`make `

# Run
build it, run the binary like this:

```bash
./koekeloere -v -w /tmp/example,/tmp/example2 http://myWebServer.tld/webhook
```

## Overview of arguments 

| Key   | Usage                                             | Example                                   |
|----   |----                                               |----                                       |
| -v    | Verbose mode (optional)                           | `-v`                                      |
| -w    | Directories to **watch**, comma seperated         | `-w /home/xantios,/tmp/example`           |
|       | all trailing arguments are parsed as URLS         | `http://MyServer.com/hook`                |

## Availbe drivers
- _http_ http and https requests (POST and GET)

## Emited events
- _create_ a file or directory is created
- _delete_ a file or directory is deleted
- _write_ change in existing file
- _chmod_ default unix permissions

# Develop a extension
see `moffel/moffel.go` for the `drivers` map. add your custom function there and add a file to `moffel/yourDriver.go` with the implementation.

the map is layed out as `map[string]interface{}`
