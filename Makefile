all:
	go build -race

test:
	go test -v

cover:
	go test -cover

coverhtml:
	go test -coverprofile=cover.out &&  go tool cover -html=cover.out

clean:
	rm cover.out macl macl.log
