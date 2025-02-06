// package main

// import (
// 	"bytes"
// 	"context"
// 	"encoding/json"
// 	"fmt"
// 	"net/http"
// 	"sync"
// 	"time"

// 	"github.com/docker/docker/api/types/container"
// 	"github.com/docker/docker/client"
// 	"github.com/go-ping/ping"
// )

// const (
// 	url = "http://localhost:8080/addresses"
// )

// type Address struct {
// 	IP                 string `json:"ip"`
// 	LastSuccessfulPing string `json:"last_successful_ping"`
// 	ResponseTime       string `json:"response_time"`
// }

// type AddressList struct {
// 	Addresses []Address `json:"addresses"`
// }

// func getAddresses() (AddressList, error) {
// 	ctx := context.Background()
// 	cli, err := client.NewClientWithOpts(client.FromEnv)
// 	if err != nil {
// 		panic(err)
// 	}

// 	containers, err := cli.ContainerList(ctx, container.ListOptions{})
// 	if err != nil {
// 		panic(err)
// 	}

// 	for _, container := range containers {
// 		containerJSON, err := cli.ContainerInspect(ctx, container.ID)
// 		if err != nil {
// 			fmt.Printf("Error inspecting container %s: %v\n", container.ID, err)
// 			continue
// 		}
// 		for networkName, network := range containerJSON.NetworkSettings.Networks {
// 			fmt.Printf("Container ID: %s, Names: %v, Network: %s, IP Address: %s\n",
// 				container.ID, container.Names, networkName, network.IPAddress)
// 		}
// 	}
// 	return AddressList{}, nil
// }

// func pingAddress(pinger *ping.Pinger) (Address, error) {

// 	pinger.Count = 1
// 	pinger.Timeout = time.Second * 1
// 	err := pinger.Run() // blocks until finished
// 	stats := pinger.Statistics()
// 	answer := Address{
// 		IP:           stats.Addr,
// 		ResponseTime: fmt.Sprintf("%v", stats.AvgRtt),
// 	}
// 	if err == nil {
// 		answer.LastSuccessfulPing = time.Now().Format(time.RFC3339)
// 	}
// 	return answer, err
// }

// func updateAddresses(address string, addressList AddressList) error {
// 	data, err := json.Marshal(addressList)
// 	if err != nil {
// 		return err
// 	}

// 	req, err := http.NewRequest(http.MethodPut, address, bytes.NewBuffer(data))
// 	if err != nil {
// 		return err
// 	}
// 	req.Header.Set("Content-Type", "application/json")

// 	client := &http.Client{}
// 	resp, err := client.Do(req)
// 	if err != nil {
// 		return err
// 	}
// 	defer resp.Body.Close()

// 	return nil
// }

// func worker(jobs <-chan *ping.Pinger, results chan<- *Address, wg *sync.WaitGroup, ctx context.Context) {
// 	// fmt.Println("worker", id, "starting")
// 	defer wg.Done()
// 	// time.Sleep(time.Second)
// 	for {
// 		select {
// 		case <-ctx.Done():
// 			// t := time.NewTimer(time.Second * time.Duration(rand.Intn(10)))
// 			// <-t.C
// 			// fmt.Println("worker", id, "cancelled")
// 			return
// 		case url, ok := <-jobs:
// 			if !ok {
// 				// fmt.Println("worker", id, "finished")
// 				return
// 			}
// 			ansver, _ := pingAddress(url)
// 			// fmt.Println("worker", id, "take", url)
// 			// time.Sleep(time.Second)
// 			results <- &ansver
// 		}
// 	}
// }

// func crawlWeb(pingerCh chan *ping.Pinger, ctx context.Context) chan *Address {
// 	const numWorker = 8
// 	results := make(chan *Address)
// 	go func() {
// 		wg := sync.WaitGroup{}
// 		for w := 1; w <= numWorker; w++ {
// 			wg.Add(1)
// 			go worker(pingerCh, results, &wg, ctx)
// 		}
// 		wg.Wait()
// 		// fmt.Println("wDone")
// 		close(results)
// 	}()

// 	return results
// }

// func main() {
// 	/*addressList, err := */ getAddresses()
// 	// if err != nil {
// 	// 	fmt.Println("Error getting addresses:", err)
// 	// 	return
// 	// }

// 	// numAddresses := len(addressList.Addresses)

// 	// pingers := make([]*ping.Pinger, numAddresses)

// 	// for i := range addressList.Addresses {
// 	// 	pingers[i], err = ping.NewPinger(addressList.Addresses[i].IP)
// 	// 	if err != nil {
// 	// 		fmt.Printf("Error create pinger%v: %v\n", i, err)
// 	// 		return
// 	// 	}
// 	// }

// 	// jobs := make(chan *ping.Pinger)
// 	// cntx, cancel := context.WithCancel(context.Background())

// 	// go func() {
// 	// 	exit := make(chan os.Signal, 1)
// 	// 	signal.Notify(exit, os.Interrupt, syscall.SIGTERM)
// 	// 	<-exit
// 	// 	cancel()
// 	// }()

// 	// go func() {
// 	// 	for j := 0; j < numAddresses; j++ {
// 	// 		jobs <- pingers[j]
// 	// 	}
// 	// 	close(jobs)
// 	// }()

// 	// results := crawlWeb(jobs, cntx)
// 	// answer := make([]Address, numAddresses)
// 	// for a := range results {
// 	// 	answer = append(answer, *a)
// 	// }

// 	// var answerList AddressList
// 	// answerList.Addresses = answer

//		// if err := updateAddresses(url, answerList); err != nil {
//		// 	fmt.Println("Error updating addresses:", err)
//		// }
//	}
package main

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

func main() {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}

	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}

	for _, container := range containers {
		containerJSON, err := cli.ContainerInspect(ctx, container.ID)
		if err != nil {
			fmt.Printf("Error inspecting container %s: %v\n", container.ID, err)
			continue
		}
		for networkName, network := range containerJSON.NetworkSettings.Networks {
			fmt.Printf("Container ID: %s, Names: %v, Network: %s, IP Address: %s\n",
				container.ID, container.Names, networkName, network.IPAddress)
		}
	}
}
