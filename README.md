# caching

caching is a Go package, implementing an in memory cache

## Installation

```bash
go get github.com/g3offrey/caching
```

## Usage

```go
import "github.com/g3offrey/caching"

func main() {
   c := caching.New[int64](100 * time.Second)

   v := c.Remember("key", expensiveComputationFunction)
   // expensiveComputationFunction is not called a second time
   v2 := c.Remember("key", expensiveComputationFunction)
}

```

You can also recompute in background and retrieve the stale value.
```go
import "github.com/g3offrey/caching"

func main() {
   c := caching.New[int64](100 * time.Second)

   v := c.GetStaleThenRecompute("key", expensiveComputationFunction)
   time.Sleep(2 * time.Hour)
   // v2 is the same as v but value is recomputed in background
   v2 := c.GetStaleThenRecompute("key", expensiveComputationFunction)
}

```

## Contributing

Pull requests are welcome. For major changes, please open an issue first
to discuss what you would like to change.

Please make sure to update tests as appropriate.

## License

[MIT](https://choosealicense.com/licenses/mit/)