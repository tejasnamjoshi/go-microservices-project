start:
	docker-compose up -d

testAll: 
	go test -timeout 30s ./...

testCov: 
	go test -coverprofile=coverage.out -timeout 30s ./...