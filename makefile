docker-build:
	docker build -t go-auth .

docker-run:
	docker run -p 8080:8080 --rm -v $(pwd):/app -v /app/tmp --name go-auth  go-auth