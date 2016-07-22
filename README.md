# redis-diff

Diff Redis Instances for key and value matching

## Installation

```
mkdir -p ~/tmp/redis-diff
cd !$
GOPATH=$(pwd) go get github.com/zph/redis-diff
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

## License
- MIT

## Prior Art
- https://github.com/adarqui/redis-transfer

## Contributions
- Contributions are welcome
- Fork it and submit PRs
