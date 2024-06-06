
all: main

main:
	go build -o nvoke . 

image:
	docker build -t nvoke .

clean: 
	rm -f nvoke
