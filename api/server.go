package api

import (
	"finder_api/crawl"
	"log"
	"net/http"
	"sync"
	"time"
)

type Server struct {
	mux   *http.ServeMux
	data  []string
	ch    chan string
	close bool
	mutex sync.Mutex
}

func NewServer() (s *Server, err error) {
	mux := http.NewServeMux()
	data := []string{}
	ch := make(chan string, 1000)
	s = &Server{mux: mux, data: data, close: false, ch: ch}
	return
}

func (s *Server) GetMux() (mux *http.ServeMux) {
	return s.mux
}

func (s *Server) Receive() {
	go func() {
		if !s.close {
			for data := range s.ch {
				items := append(s.data, data)
				locked := s.mutex.TryLock()
				if locked {
					s.data = items
					//fmt.Println(s.data)
					s.mutex.Unlock()
				}

			}
		}

	}()
}

func (s *Server) Timer() {
	go func() {
		time.Sleep(30 * time.Second)
		close(s.ch)
		log.Println("data: ", s.data)
		s.close = true
	}()
}

func (s *Server) Crawl(url string) {

	go func() {
		for !s.close {
			results, err := crawl.FindUrls(url, s.data)
			if err != nil {
				log.Printf("error at crawl %v", err)
				continue
			}

			for _, result := range results {
				if s.close {
					break
				}
				s.ch <- result
			}
		}
	}()

}

func (s *Server) SetupRouter() {
	s.mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		value, err := w.Write([]byte("Created user"))
		if err != nil {
			log.Println("Error")
		} else {
			log.Println(":", value)
		}
	})
	s.mux.HandleFunc("GET /go", func(w http.ResponseWriter, r *http.Request) {
		value, err := w.Write([]byte("Going user"))
		if err != nil {
			log.Println("Error")
		} else {
			log.Println(": ", value)
		}
	})
}
