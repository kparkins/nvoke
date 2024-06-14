
all: main

main:
	go build -o bin/nvoke . 

image:
	docker build -t nvoke .

clean: 
	rm -f bin/nvoke
