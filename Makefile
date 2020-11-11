BIN=GrueBot


build:
	go build -o ${BIN}

build-debug:
	go build -o ${BIN} -gcflags="-l -N"

clean:
	rm ${BIN}