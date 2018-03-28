Siknäs-skylt Thumbnail Generator
================================

```bash
docker build -t siknas-skylt-thumbgen-go .

# Get dependencies
docker run -it --rm -v $(pwd):/go/src/app siknas-skylt-thumbgen-go dep ensure

docker run -it --rm -v $(pwd):/go/src/app siknas-skylt-thumbgen-go
```

