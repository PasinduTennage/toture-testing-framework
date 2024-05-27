package dummy

import (
	"encoding/json"
	"github.com/rs/cors"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type Metrics struct {
	Throughput float64 `json:"throughput"`
	Latency    float64 `json:"latency"`
	CPUUsage   float64 `json:"cpuUsage"`
	MemUsage   float64 `json:"memUsage"`
}

var (
	metrics Metrics
	mu      sync.Mutex
)

func ListenFrontEnd(name string) {

	mux := http.NewServeMux()
	mux.HandleFunc("/metrics", handleMetrics)

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"}, // Allow requests from any origin
		AllowedMethods: []string{"GET"},
	})

	handler := c.Handler(mux)

	go generateMetrics()

	log.Println("Server starting on :" + name)
	log.Fatal(http.ListenAndServe(":"+name, handler))

}

func handleMetrics(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()
	//fmt.Printf("%v\n", r.Header)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(metrics)
}

func generateMetrics() {
	for {
		mu.Lock()
		metrics = Metrics{
			Throughput: rand.Float64() * 100,
			Latency:    rand.Float64() * 100,
			CPUUsage:   rand.Float64() * 100,
			MemUsage:   rand.Float64() * 100,
		}
		mu.Unlock()
		time.Sleep(1 * time.Second)
	}
}

func DoUi(pr *Proxy) {
	ListenFrontEnd(strconv.FormatInt(pr.name*10000+100, 10))
}
