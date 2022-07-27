# `Cars fleet service` 
* This is an example service that will demonstrate how to read a json file and service the results via API endpoints, 

Run server:
go run main.go

Run test:
go test ./... -v

Build for container:
go build -o bin/cars-fleet-service main.go
docker build -f Dockerfile -t cars-fleet-service:0.0.1 .

Task: 
Build and deploy a web service to return data (no create,update,delete): 
 
Return the information of the car based on an exact value match of the “Name” key in 
the json array via the web service.

Return the information of one or multiple cars based on a search value of the “Name” 
key in the json array via the web service.

Search based APIs lend themselves to a search engine storage layer for quick retrival of data.

Potential improvements could include support for multiple json files to read and maybe support for different kinds of files such as csv. These kinds of decisions depend on project reqirements, scope etc