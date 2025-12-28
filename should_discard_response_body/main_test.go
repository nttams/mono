package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/http/httptrace"
	"testing"
)

/*
go test -bench . -benchmem
goos: darwin
goarch: arm64
pkg: should_discard_response_body
cpu: Apple M2
BenchmarkHTTPDiscard-8             25954             45770 ns/op           38605 B/op         69 allocs/op
--- BENCH: BenchmarkHTTPDiscard-8
    main_test.go:75: Connections Created: 1, Connections Reused: 0
    main_test.go:75: Connections Created: 1, Connections Reused: 99
    main_test.go:75: Connections Created: 1, Connections Reused: 9999
    main_test.go:75: Connections Created: 1, Connections Reused: 25953
BenchmarkHTTPNoDiscard-8           10000            105554 ns/op           51082 B/op        131 allocs/op
--- BENCH: BenchmarkHTTPNoDiscard-8
    main_test.go:75: Connections Created: 1, Connections Reused: 0
    main_test.go:75: Connections Created: 100, Connections Reused: 0
    main_test.go:75: Connections Created: 10000, Connections Reused: 0
PASS
ok      should_discard_response_body    3.034s
*/

func BenchmarkHTTPDiscard(b *testing.B) {
	doBench(b, true)
}

func BenchmarkHTTPNoDiscard(b *testing.B) {
	doBench(b, false)
}

func doBench(b *testing.B, discard bool) {
	server := setupServer()
	defer server.Close()

	client := &http.Client{}

	var created, reused int

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		trace := &httptrace.ClientTrace{
			GotConn: func(connInfo httptrace.GotConnInfo) {
				if connInfo.Reused {
					reused++
				} else {
					created++
				}
			},
		}

		req, _ := http.NewRequest("GET", server.URL, nil)
		req = req.WithContext(httptrace.WithClientTrace(req.Context(), trace))

		resp, err := client.Do(req)
		if err != nil {
			b.Fatal(err)
		}

		if discard {
			io.Copy(io.Discard, resp.Body)
		}
		resp.Body.Close()
	}

	// Print results at the end of the benchmark
	b.Logf("Connections Created: %d, Connections Reused: %d", created, reused)
}

func setupServer() *httptest.Server {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(make([]byte, 32*1024))
	})
	return httptest.NewServer(handler)
}
