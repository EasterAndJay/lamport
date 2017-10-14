CLIENT_COUNT=5


run:
	go build main.go
	for number in 0 1 2 3 4
	./main -n ${CLIENT_COUNT} -pid