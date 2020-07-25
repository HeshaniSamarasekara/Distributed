# Distributed Project P2P

## How to setup
#### Prerequisites
- Go Installed (go version go1.12.6 or higher)
- JDK installed (13.0.2 higher)
- Linux nodes (preferably CentOS 7 or higher)

**Step1:** Copy the Project to your $GOPATH

**Step 2:** Navigate into the P2PClient directory.
```bash
cd $GOPATH/src/Distributed/P2PClient/
```
**Step 3:** Build the Project
```bash
go build
```

## How to start a node

### On Linux/iOS

```
./P2PClient [IP] [PORT] [NODE NAME]
```
Ex:
```
./P2PClient 127.0.0.1 9000 NODE1
```

### On Windows

```
.\P2PClient.exe [IP] [PORT] [NODE NAME]
```
Ex:
```
.\P2PClient.exe 127.0.0.1 9000 NODE1
```

## How to register a node in the system

This command is not necessary to run. When a node starts, it automatically sends the register call to the Bootstrap server.
Ex:
```
POST http://127.0.0.1:9000/register
```

## How to unregister a node in the system
Ex:
```
DELETE http://127.0.0.1:9000/unregister
```

## How to search files of a node
Ex:
```
GET http://127.0.0.1:9000/files
```

## How to search routing table of a node
Ex:
```
GET http://127.0.0.1:9000/routeTable
```

## How to search a file in the system
Ex:
```
GET http://127.0.0.1:9000/search/{file_name}
```
Expected sample response
```
[SIZE]  SEROK [NO_OF_FILES] [HOST]  [PORT]  [HOP_COUNT] [FILE_NAME]

0035 SEROK 1 localhost 1111 0 Glee
```
For filenames with more than one word should replace spaces in the name with underscore character.
Ex: Harry_Potter

## How to search File Table of a node
Ex:
```
GET http://127.0.0.1:9000/fileTable
```

## How to download File from a node
Replace the parameters in the below command with the search result of the file needed.
Ex:
```
GET http://localhost:9000/download/{server}/{port}/{file_name}
```

## Tested configurations
- Java Bootstrap Server
- CentOS Linux release 7.7.1908 (Core)
- openjdk version "13.0.2" 2020-01-14
- go version go1.12.6 linux/amd64
