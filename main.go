package main

import (
	"fmt"
	"net"
	"sort"
)

// func main() {
// 	for i := 21; i < 120; i++ {
// 		address := fmt.Sprintf("127.0.0.1:%d", i)
// 		conn, err := net.Dial("tcp", address)
// 		if err != nil {
// 			fmt.Printf("%s closed\n", address)
// 			continue
// 		}
// 		conn.Close()
// 		fmt.Printf("%s Opened\n", address)
// 	}
// }

// func main() {
// 	start := time.Now()
// 	var wg sync.WaitGroup
// 	for i := 1; i < 65535; i++ {
// 		wg.Add(1)
// 		go func(j int) {
// 			defer wg.Done()
// 			address := fmt.Sprintf("127.0.0.1:%d", j)
// 			conn, err := net.Dial("tcp", address)
// 			if err != nil {
// 				fmt.Printf("%s closed\n", address)
// 				return
// 			}
// 			conn.Close()
// 			fmt.Printf("%s Opened\n", address)
// 		}(i)
// 	}
// 	wg.Wait()
// 	elapsed := time.Since(start)/1e9
// 	fmt.Printf("/n/n%d seconds", elapsed)
// }

func worker(ports chan int, results chan int) {
	for p := range ports {
		address := fmt.Sprintf("127.0.0.1:%d", p)
		conn, err := net.Dial("tcp", address)
		if err != nil {
			fmt.Printf("Error: %s\n", address)
			results <- 0
			continue
		}
		conn.Close()
		fmt.Printf("%s opened\n", address)
		results <- p
	}
}

func main() {
	ports := make(chan int, 1024)
	results := make(chan int)
	var openports []int
	var closeports []int

	for i := 0; i < cap(ports); i++ {
		go worker(ports, results)
	}

	go func() {
		for i := 0; i < 1024; i++ {
			ports <- i
		}
	}()

	for i := 0; i < 1024; i++ {
		port := <-results
		if port != 0 {
			openports = append(openports, port)
		} else {
			closeports = append(closeports, port)
		}
	}

	close(ports)
	close(results)

	sort.Ints(openports)
	sort.Ints(closeports)
	for _, port := range openports {
		fmt.Printf("%d opened\n", port)
	}
	for _, port := range closeports {
		fmt.Printf("%d closed\n", port)
	}

}
