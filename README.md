# Clipshare

A clipboard sharing tool written in Go.
The tool currently works for linux distros with X Windows and requires xsel installed.

## Usage

### Start the tool
```
clipshare start
```
Ideally run as a background process.

### Get contents of clipshare buffer (received from other peers)
```
clipshare get
```

### Send clipboard contents to other peers
```
clipshare set 127.0.0.1 10.0.2.1
```

### Stop clipshare
```
clipshare stop
```


###Todo
* Queued Clipshare buffer
* Daemon mode
* Security
* refactor
  
