# Distributed
Distributed Project P2P

to build use [go build]

How to Run a node

For windows
.\P2PClient.exe [IP] [PORT] [NODE NAME]
Ex:
.\P2PClient.exe 127.0.0.1 9000

How to register a node in the system
Ex:
 POST http://127.0.0.1:9000/register

How to unregister a node in the system
Ex: 
DELETE http://127.0.0.1:9000/unregister

How to search files of a node
Ex: 
GET http://127.0.0.1:9000/files

How to search routing table of a node
Ex: 
GET http://127.0.0.1:9000/routeTable

How to search a file in the system
Ex: 
GET http://127.0.0.1:9000/search/{file_name}

How to search file Table of a node
Ex: 
GET http://127.0.0.1:9000/fileTable
