# Chat Service
This backend service handles the persistent websockets between users as well as other chat related endpoints. It is written in go and uses Gin and go-socket-io.

## Dependencies

Please ensure that you have atleast go1.15.6 and you are running this on a Unix based system. 

## Running the service

Please view the main CATchat repository to deploy the 
entire service. To run only the Chat service, run the following command.

```
go run cmd/main.go
```
This can also be run through Docker. 
```
docker run .
```

## Testing this service
Please ensure that you have Go installed on your system.
Running command in the root directory will run the tests and display the test coverage.

```
go test ./... -cover
```

To view the actual code coverage, run the commands below.
```
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```
