Siknäs-skylt Thumbnail Generator
================================

```bash
docker build -t siknas-skylt-thumbgen-go .

# Get dependencies
docker run -it --rm -v $(pwd):/go/src/app siknas-skylt-thumbgen-go dep ensure -v

# Help
docker run -it --rm -v $(pwd):/go/src/app siknas-skylt-thumbgen-go go run main.go --help

# Connect to the running server
docker run -it --rm -v $(pwd):/go/src/app siknas-skylt-thumbgen-go go run *.go --host $(docker-machine ip):8080
```
