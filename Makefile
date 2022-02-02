## test: runs all tests
test:
	@go test -v ./...

## cover: opens coverage in browser
cover:
	@go test -coverprofile=coverage.out ./... && go tool cover -html=coverage.out

## coverage: displays test coverage
coverage:
	@go test -cover ./...

## build: builds the api
build:
	cd api && go build -o api

## mock: creates mock files
mock:
	cd api && mockery --all

## up: docker compose up
up:
	docker-compose -f docker-compose.yaml up -d  --build --remove-orphans 

## up: docker compose down
down:
	docker-compose -f docker-compose.yaml down 

## mongo-req: successful mongo req
mongo-req:
	curl -X POST \
			-H "Content-type: application/json" \
			-H "Accept: application/json" \
			-d '{"startDate": "2016-01-26","endDate": "2018-02-02","minCount": 2700,"maxCount": 3000}' \
			"http://localhost:3001/mongo"

mongo-req2:
	curl -X POST \
			-H "Content-type: application/json" \
			-H "Accept: application/json" \
			-d '{"startDate": "2016-01-26","endDate": "2018-22-02","minCount": 2700,"maxCount": 3000}' \
			"http://localhost:3001/mongo"

inmemory-post:
	curl -X POST \
			-H "Content-type: application/json" \
			-H "Accept: application/json" \
			-d '{"key": "active-tabs","value": "getir"}' \
			"http://localhost:3001/in-memory"

inmemory-post2:
	curl -X POST \
			-H "Content-type: application/json" \
			-H "Accept: application/json" \
			-d '{"value": "getir"}' \
			"http://localhost:3001/in-memory"

inmemory-get:
	curl -X GET \
			-H "Content-type: application/json" \
			-H "Accept: application/json" \
			"http://localhost:3001/in-memory?key=active-tabs"

## apply: spins up an ec2 instance and install go files
apply:
	rm -f api/api && cd infra &&  terraform apply -auto-approve

## apply: destroy ec2 instance
destroy:
	cd infra && terraform destroy -auto-approve