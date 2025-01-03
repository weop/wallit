.PHONY: build clean

build:
	go get github.com/golang/freetype
	go build -o wallit

clean:
	rm -f wallit