package infoServer

import (
	"fmt"
	"net/http"
	"runtime"
)

func HandleGetInfoServer(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Операционная система:", runtime.GOOS)
	fmt.Fprintln(w, "Количество ядер процессора:", runtime.NumCPU())
	fmt.Fprintln(w, "Количество горутин:", runtime.NumGoroutine())

	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	bToMb := func(b uint64) float64 { return float64(b) / 1024 / 1024 }

	fmt.Fprintf(w, "Alloc = %.2f MiB\n", bToMb(m.Alloc))
	fmt.Fprintf(w, "TotalAlloc = %.2f MiB\n", bToMb(m.TotalAlloc))
	fmt.Fprintf(w, "Sys = %.2f MiB\n", bToMb(m.Sys))
	fmt.Fprintf(w, "HeapAlloc = %.2f MiB\n", bToMb(m.HeapAlloc))
	fmt.Fprintf(w, "HeapSys = %.2f MiB\n", bToMb(m.HeapSys))
	fmt.Fprintf(w, "StackInuse = %.2f MiB\n", bToMb(m.StackInuse))
	fmt.Fprintf(w, "StackSys = %.2f MiB\n", bToMb(m.StackSys))
}
