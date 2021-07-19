package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"runtime"
	"time"

	"debugging/profiling/upload"

	"github.com/cheggaaa/pb/v3"
	"golang.org/x/sync/errgroup"
)

func main() {
	// limit number of CPU that program can use.
	runtime.GOMAXPROCS(runtime.NumCPU() / 3)

	start := time.Now()

	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	ctx := context.Background()
	eg, ctx := errgroup.WithContext(ctx)

	attempts := 100000
	p := pb.StartNew(attempts)

	concurrency := 100
	guard := make(chan struct{}, concurrency)

	for i := 0; i < attempts; i++ {
		guard <- struct{}{}

		eg.Go(func() error {
			defer func() { <-guard }()

			f, err := os.Open("profiling/upload/file.txt")
			if err != nil {
				return fmt.Errorf("failed to open: %w", err)
			}

			if err := upload.Do(ctx, f, upload.NotFixed); err != nil {
				return fmt.Errorf("failed to upload: %w", err)
			}

			p.Increment()
			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		fmt.Println("failed to complete the test:", err)
		return
	}

	p.Finish()
	fmt.Println("test program completed in", time.Since(start))

	fmt.Println("sleep 3 minutes to collect results")
	time.Sleep(time.Minute * 3)
}
