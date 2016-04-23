all:
	go build -o bin/bench ./bench
	go build -o bin/chat-original ./chat-original
	go build -o bin/chat-step1 ./chat-step1
	go build -o bin/chat-step2 ./chat-step2
