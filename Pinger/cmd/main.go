package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
)

const (
	url = "http://backend:8081/put_address"
)

type Address struct {
	IP                 string    `json:"ip"`
	LastSuccessfulPing time.Time `json:"last_successful_ping"`
	ResponseTime       string    `json:"response_time"`
}

type AddressList struct {
	Addresses []Address `json:"addresses"`
}

func pingAddress(address string) error {
	var answer []Address
	ping := time.Now()
	conn, err := net.DialTimeout("tcp", address, 2*time.Second)
	duration := time.Since(ping)
	if err != nil {
		fmt.Printf("%s %d (timeout) %v\n", address, duration.Milliseconds(), err)
		answer = []Address{{
			ResponseTime: "not answer",
		}}
	} else {
		conn.Close()
		answer = []Address{{
			ResponseTime:       fmt.Sprintf("%d ms", duration.Milliseconds()),
			LastSuccessfulPing: time.Now(),
		}}
	}

	answer[0].IP = strings.Split(address, ":")[0]

	var answerList AddressList
	answerList.Addresses = answer

	if err := updateAddresses(url, answerList); err != nil {
		fmt.Println("Error updating addresses:", err)
	}

	fmt.Printf("%s %d %v\n", address, duration.Milliseconds(), time.Now().Format(time.RFC3339))
	return nil
}

func updateAddresses(address string, addressList AddressList) error {
	data, err := json.Marshal(addressList)
	if err != nil {
		return err
	}

	resp, err := http.Post(address, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	return nil
}

func worker(id int, jobs <-chan string, wg *sync.WaitGroup, ctx context.Context) {
	fmt.Println("worker", id, "starting")
	defer wg.Done()
	// time.Sleep(time.Second)
	for {
		select {
		case <-ctx.Done():
			// t := time.NewTimer(time.Second * time.Duration(rand.Intn(10)))
			// <-t.C
			fmt.Println("worker", id, "cancelled")
			return
		case url, ok := <-jobs:
			if !ok {
				fmt.Println("worker", id, "finished")
				return
			}
			fmt.Println("worker", id, "take", url)
			err := pingAddress(url)
			if err != nil {
				fmt.Println("Error pinging:", err)
			}
			// time.Sleep(time.Second)
		}
	}
}

func crawlWeb(pingerCh chan string, ctx context.Context) {
	const numWorker = 3
	go func() {
		wg := sync.WaitGroup{}
		for w := 1; w <= numWorker; w++ {
			wg.Add(1)
			go worker(w, pingerCh, &wg, ctx)
		}
		wg.Wait()
		// fmt.Println("wDone")
	}()
}

func main() {

	cntx, cancel := context.WithCancel(context.Background())
	jobs := make(chan string)

	go func() {
		exit := make(chan os.Signal, 1)
		signal.Notify(exit, os.Interrupt, syscall.SIGTERM)
		<-exit
		close(jobs)
		cancel()
	}()

	crawlWeb(jobs, cntx)
	ticker := time.NewTicker(20 * time.Second)
	fmt.Println("ticker!")
	defer ticker.Stop()
	for range ticker.C {

		cli, err := client.NewClientWithOpts(
			client.FromEnv,
			client.WithVersion("1.41"),
		)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("%v\n", cli)

		containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{
			Filters: filters.NewArgs(filters.Arg("status", "running")),
		})
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Printf("faund %v containers\n", len(containers))
		for _, c := range containers {

			fmt.Printf("f %v container\n", c.Names)
			// Получаем информацию о контейнере
			inspect, err := cli.ContainerInspect(context.Background(), c.ID)
			if err != nil {
				fmt.Printf("Error inspecting %s: %v\n", c.ID[:12], err)
				continue
			}
			// Получаем все IP-адреса контейнера
			var ips []string
			for _, net := range inspect.NetworkSettings.Networks {
				ips = append(ips, net.IPAddress)
			}
			numAddresses := len(ips)
			if numAddresses == 0 {
				fmt.Printf("%s: No IP addresses found\n", c.Names[0])
				continue
			}

			fmt.Printf("f %v ip\n", ips)

			// Получаем первый экспознутый порт
			exposedPorts := inspect.Config.ExposedPorts
			if len(exposedPorts) == 0 {
				fmt.Printf("%s: No exposed ports\n", c.Names[0])
				continue
			}

			fmt.Printf("f %v exposedPorts\n", exposedPorts)

			var firstPort string
			for port := range exposedPorts {
				firstPort = port.Port()
				break
			}

			fmt.Printf("f %v firstPort\n", firstPort)

			go func(addresses []string) {
				for j := range addresses {
					jobs <- fmt.Sprintf("%s:%s", addresses[j], firstPort)
				}
			}(ips)

		}

	}
}
