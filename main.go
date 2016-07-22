package main

import (
	"flag"
	"fmt"
	"gopkg.in/redis.v4"
	"log"
	"net"
	"net/url"
	"os"
	"reflect"
	"strconv"
	"sync"
)

func parseURI(s string) (server *Server, err error) {
	// Defaults
	host := "localhost"
	password := ""
	port := 6379
	db := 0

	u, err := url.Parse(s)
	if err != nil {
		log.Fatal(err)
	}
	if u.Scheme != "redis" {
		log.Fatal("Scheme must be redis")
	}
	q := u.Query()
	dbS := q.Get("db")
	if u.User != nil {
		var ok bool
		password, ok = u.User.Password()
		if !ok {
			password = ""
		}
	}

	var p string
	host, p, _ = net.SplitHostPort(u.Host)

	if p != "" {
		port, err = strconv.Atoi(p)
		if err != nil {
			log.Fatalf("Unable to convert port to integer for %s", err)
		}
	}

	if dbS != "" {
		db, err = strconv.Atoi(dbS)
		if err != nil {
			log.Fatalf("Unable to convert db to integer for %s", dbS)
		}
	}

	client := CreateClient(host, password, port, db)
	return &Server{client, host, port, db, password}, nil
}

func (s *Server) scanner(match pattern, wg *sync.WaitGroup) chan redisKey {
	keyChan := make(chan redisKey, 1000)
	split := make(chan []string)

	splitter := func() {
		defer wg.Done()
		defer close(keyChan)
		for {
			select {
			case ks, ok := <-split:
				if !ok {
					return
				}
				for _, k := range ks {
					keyChan <- redisKey(k)
				}
			}
		}
	}

	keyScanner := func(c chan redisKey) {
		defer wg.Done()
		var cursor uint64
		var n int
		for {
			var keys []string
			var err error
			// REDIS SCAN
			// http://redis.io/commands/scan
			// Preferable because it doesn't lock complete database on larger keysets for 250ms+.
			keys, cursor, err = s.client.Scan(cursor, string(match), 1000).Result()
			if err != nil {
				log.Fatal("KeysRedis: error obtaining keys list from redis: ", err)
			}
			split <- keys

			n += len(keys)
			if cursor == 0 {
				close(split)
				return
			}
		}
	}

	wg.Add(1)
	go splitter()

	wg.Add(1)
	go keyScanner(keyChan)

	return keyChan
}

func (p *Pipe) compare(src, dst *Server, key redisKey) (interface{}, interface{}, bool) {
	s, err := src.client.Get(string(key)).Result()
	if err != nil {
		log.Printf("Unable to get expected key %s from src: %+v", key, src.client)
	}
	d, _ := dst.client.Get(string(key)).Result()
	isMatch := reflect.DeepEqual(s, d)
	return s, d, isMatch

}

func (p *Pipe) CompareKeys(c chan redisKey, mismatches chan *Discrepancy, wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case _, ok := <-p.shutdown:
				if !ok {
					return
				}
			case k, ok := <-c:
				if !ok {
					return
				}
				s, d, isMatch := p.compare(p.from, p.to, k)
				if !isMatch {
					mismatches <- &Discrepancy{k, s, d}
				}
			}
		}
	}()
}

func CreateClient(host, pass string, port, db int) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", host, port),
		Password: pass,
		DB:       db,
	})
}

func writer(c chan *Discrepancy, wg *sync.WaitGroup, del *string) {
	defer wg.Done()
	i := *del
	for d := range c {
		fmt.Printf("%s%s%s%s%s\n", d.key, i, d.src, i, d.dst)
	}
}

func main() {
	src := flag.String("src", "", "Format redis://:password@host:port?db=0")
	dst := flag.String("dst", "redis://localhost:6379", "redis://:password@host:port?db=0")
	threads := flag.Int("parallel", 20, "Threading count. Default `20`")
	match := flag.String("keys", "*", "Match subset of keys `*`")
	delimiter := flag.String("delimiter", "|", "Delimiter that will be used to separate output")
	flag.Parse()
	if *src == "" {
		flag.Usage()
		os.Exit(1)
	}
	from, _ := parseURI(*src)
	to, _ := parseURI(*dst)

	var wg sync.WaitGroup
	shutdown := make(chan bool, 1)
	discrepancies := make(chan *Discrepancy)
	pipe := &Pipe{from, to, *match, shutdown}
	keyChan := pipe.from.scanner(pattern(*match), &wg)

	tx := *threads
	for i := 0; i < tx; i++ {
		p := &Pipe{from, to, *match, shutdown}
		p.CompareKeys(keyChan, discrepancies, &wg)
	}

	// Setup Writer
	var wgWriter sync.WaitGroup
	wgWriter.Add(1)
	go writer(discrepancies, &wgWriter, delimiter)

	// Wait for threads to complete
	wg.Wait()
	// Start cleanup routine for writer
	close(discrepancies)
	// Wait for writer to close fn
	wgWriter.Wait()
}
