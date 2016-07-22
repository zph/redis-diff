# redis-diff

Diff Redis Instances for key and value matching. Uses incremental REDIS SCAN
with cursor to limit the load impact on REDIS.

## Installation

```
mkdir -p ~/tmp/redis-diff
cd !$
GOPATH=$(pwd) go get github.com/zph/redis-diff
cp bin/redis-diff ~/bin/ # [or somewhere else on your path]
```

## Usage

``` bash
$ ./bin/redis-diff
Usage of ./bin/redis-diff:
  -delimiter string
      Delimiter that will be used to separate output (default "|")
  -dst string
      redis://:password@host:port?db=0 (default "redis://localhost:6379")
  -keys *
      Match subset of keys * (default "*")
  -parallel 20
      Threading count. Default 20 (default 20)
  -src string
      Format redis://:password@host:port?db=0
```

Sample Output on STDOUT

```
key1|valuedb1|valuedb2
key4|valuedb1|valuedb2
key6|valuedb1|valuedb2

```

## Performance
Iterates through src keyspace at ~500 to 1000 keys per second on my development
machine when pointed at hosted redis instance on separate machine.
I suspect the bottleneck is the key scanner itself, but I didn't see ways
to improve that performance.

## License
- MIT

## Prior Art
- https://github.com/adarqui/redis-transfer

## Contributions
- Contributions are welcome
- Fork it and submit PRs
