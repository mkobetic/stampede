package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"runtime/pprof"
	"runtime/trace"
	"strings"
	"sync"
)

var (
	root   = flag.String("root", "$HOME/st/mail", "directory containing the mail archive")
	assets = flag.String("assets", "assets", "directory containing asset files (css, ...)")
	// profiling options
	cpuf   = flag.String("cpu", "", "filename for a CPU profile")
	memf   = flag.String("mem", "", "filename for a heap profile")
	tracef = flag.String("trace", "", "filename for an event trace")
	// testing options
	dump = flag.Bool("dump", false, "dump debugging details for all folders and quit")

	Root *MailDirectory
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	flag.Parse()

	if *cpuf != "" {
		if f, err := os.Create(*cpuf); err != nil {
			log.Fatalf("cpu profile: ", err.Error())
		} else {
			if err := pprof.StartCPUProfile(f); err != nil {
				log.Fatalf("start cpu profile: ", err.Error())
			} else {
				defer func(f *os.File) {
					pprof.StopCPUProfile()
					f.Close()
				}(f)
			}
		}
	}

	if *memf != "" {
		if f, err := os.Create(*memf); err != nil {
			log.Fatalf("mem profile: ", err.Error())
		} else {
			defer func(f *os.File) {
				pprof.Lookup("heap").WriteTo(f, 0)
				f.Close()
			}(f)
		}
	}

	if *tracef != "" {
		if f, err := os.Create(*tracef); err != nil {
			log.Fatalf("trace: ", err.Error())
		} else {
			if err := trace.Start(f); err != nil {
				log.Fatalf("start trace: ", err.Error())
			} else {
				defer func(f *os.File) {
					trace.Stop()
					f.Close()
				}(f)
			}
		}
	}

	var wg sync.WaitGroup
	Root = OpenDirectory(os.ExpandEnv(*root), nil, &wg)
	wg.Wait()

	if *dump {
		Root.DumpOffsets()
		os.Exit(0)
	}

	http.HandleFunc("/", Navigate)
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir(*assets))))
	go func() {
		log.Fatal(http.ListenAndServe(":8080", nil))
	}()
	log.Print("Listening at http://localhost:8080")

	// handle INT so that profiles and traces can flush properly
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

}

func Navigate(w http.ResponseWriter, r *http.Request) {
	var path []string
	// URL.Path always starts with a /
	path = strings.Split(r.URL.Path[1:], "/")
	if len(path) > 0 && len(path[0]) == 0 {
		path = path[1:]
	}
	log.Printf("Navigate %#v", path)
	if d := Root.Find(path); d != nil {
		d.ServeHTTP(w, r)
	} else {
		http.Error(w, "Not Found", 404)
	}
}
