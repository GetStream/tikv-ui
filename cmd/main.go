package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/GetStream/tikv-ui/pkg/handlers"
	"github.com/GetStream/tikv-ui/pkg/server"
	"github.com/GetStream/tikv-ui/pkg/utils"
	"github.com/tikv/client-go/v2/config"
	"github.com/tikv/client-go/v2/rawkv"
)

func main() {
	// Read PD addresses from env: TIKV_PD_ADDRS="127.0.0.1:2379,127.0.0.1:2381"
	pdAddrsEnv := os.Getenv("TIKV_PD_ADDRS")
	if pdAddrsEnv == "" {
		log.Fatal("TIKV_PD_ADDRS env var is required (comma-separated PD addresses)")
	}
	clusterStrs := strings.Split(strings.Trim(pdAddrsEnv, ";"), ";")
	pdAddrs := utils.SplitAndTrim(clusterStrs[0], ",")

	ctx := context.Background()

	// Create TiKV RawKV client for default cluster
	cli, err := rawkv.NewClient(ctx, pdAddrs, config.DefaultConfig().Security)
	if err != nil {
		log.Fatalf("failed to create TiKV RawKV client: %v", err)
	}

	log.Printf("Connected to default TiKV cluster ID: %d", cli.ClusterID())

	srv := server.New(cli, pdAddrs)
	defer srv.Close()

	for i, cluster := range clusterStrs[1:] {
		srv.AddCluster(ctx, fmt.Sprintf("cluster-%d", time.Now().Unix()+int64(i)), utils.SplitAndTrim(cluster, ","))
	}

	mux := http.NewServeMux()

	// Health check
	mux.HandleFunc("/health", handlers.Health(srv))

	// Cluster management
	mux.HandleFunc("/api/clusters", handlers.ListClusters(srv))
	mux.HandleFunc("/api/clusters/connect", handlers.Connect(srv))
	mux.HandleFunc("/api/clusters/switch", handlers.SwitchCluster(srv))

	// Raw KV operations
	mux.HandleFunc("/api/raw/get", handlers.Get(srv))
	mux.HandleFunc("/api/raw/put", handlers.Put(srv))
	mux.HandleFunc("/api/raw/delete", handlers.Delete(srv))
	mux.HandleFunc("/api/raw/scan", handlers.Scan(srv))

	// Serve static files (Frontend)
	// In Docker, we will copy the built 'out' directory to 'public'
	staticDir := "./public"
	fs := http.FileServer(http.Dir(staticDir))
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Calculate the file path
		path := staticDir + r.URL.Path

		// Check if file exists, otherwise serve index.html (SPA fallback)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			http.ServeFile(w, r, staticDir+"/index.html")
			return
		}

		fs.ServeHTTP(w, r)
	})
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}
	httpServer := &http.Server{
		Addr:         ":" + port,
		Handler:      server.CORSMiddleware(server.LoggingMiddleware(mux)),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// Graceful shutdown
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
		<-sigChan

		log.Println("Shutting down server...")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			log.Printf("Server shutdown error: %v", err)
		}
	}()

	log.Println("TiKV explorer API listening on :" + port)
	log.Println("Cluster management: POST /api/clusters/connect, GET /api/clusters, POST /api/clusters/switch")
	if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server error: %v", err)
	}
}
