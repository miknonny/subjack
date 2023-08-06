package subjack

import (
	"fmt"
	"log"
	"sync"
)

type Options struct {
	Domain       string
	Wordlist     string
	Threads      int
	Timeout      int
	Output       string
	Ssl          bool
	All          bool
	Verbose      bool
	Config       string
	Manual       bool
	Fingerprints []Fingerprints
}

type Subdomain struct {
	Url string
}

/* Start processing subjack from the defined options. */
func Process(o *Options) {
	var list []string
	var err error

	urls := make(chan *Subdomain, o.Threads)
	fmt.Println(o.Domain)
	if len(o.Domain) > 0 {
		list = append(list, o.Domain)
	} else {
		list, err = open(o.Wordlist)
	}

	if err != nil {
		log.Fatalln(err)
	}

	o.Fingerprints = fingerprints(o.Config)

	// Setting up workers.
	var wg sync.WaitGroup
	wg.Add(o.Threads)
	for i := 0; i < o.Threads; i++ {

		go func() {
			defer wg.Done()
			for url := range urls {
				url.dns(o)
			}
		}()
	}

	for i := 0; i < len(list); i++ {
		urls <- &Subdomain{Url: list[i]}
	}

	close(urls)
	wg.Wait()
}
