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

	"main/internal/config"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/sirupsen/logrus"
)

type Address struct {
	IP                 string    `json:"ip"`
	LastSuccessfulPing time.Time `json:"last_successful_ping"`
	ResponseTime       string    `json:"response_time"`
}

type AddressList struct {
	Addresses []Address `json:"addresses"`
}

var cnfg *config.Config

func pingAddress(address string) error {
	var answer []Address
	ping := time.Now()
	conn, err := net.DialTimeout("tcp", address, 2*time.Second)
	duration := time.Since(ping)
	if err != nil {
		logrus.Warnf("%s %d (timeout) %v\n", address, duration.Milliseconds(), err)
		answer = []Address{{
			ResponseTime: "no answer",
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

	if err := updateAddresses(cnfg.Address, answerList); err != nil {
		logrus.Errorln("Error updating addresses:", err)
	}

	logrus.Debugf("%s %d %v\n", address, duration.Milliseconds(), time.Now().Format(time.RFC3339))
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
	logrus.Infoln("worker", id, "starting")
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():

			logrus.Infoln("worker", id, "cancelled")
			return
		case url, ok := <-jobs:
			if !ok {
				logrus.Infoln("worker", id, "finished")
				return
			}
			logrus.Debugln("worker", id, "take", url)
			err := pingAddress(url)
			if err != nil {
				logrus.Errorln("Error pinging:", err)
			}
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
	}()
}

func main() {

	cnfg = config.ParseConfig("config/config.yaml")

	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	logrus.SetLevel(logrus.Level(cnfg.Level))

	cntx, cancel := context.WithCancel(context.Background())
	jobs := make(chan string)

	logrus.Debugf("Configuration loaded: %+v\n", cntx)

	go func() {
		exit := make(chan os.Signal, 1)
		signal.Notify(exit, os.Interrupt, syscall.SIGTERM)
		<-exit
		close(jobs)
		cancel()
	}()

	crawlWeb(jobs, cntx)
	ticker := time.NewTicker(20 * time.Second)

	logrus.Debugln("ticker!")

	defer ticker.Stop()
	for range ticker.C {

		cli, err := client.NewClientWithOpts(
			client.FromEnv,
			client.WithVersion("1.41"),
		)
		if err != nil {
			logrus.Errorln(err)
			return
		}

		logrus.Debugf("Client: %+v\n", cli)

		containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{
			Filters: filters.NewArgs(filters.Arg("status", "running")),
		})
		if err != nil {
			logrus.Errorln(err)
			return
		}

		logrus.Debugf("faund %v containers\n", len(containers))

		for _, c := range containers {

			logrus.Debugf("f %v container\n", c.Names)

			if strings.Index(c.Names[0], "docker-pinger") != -1 {
				logrus.Debugln("oh... it`s me!?")
				continue
			}

			inspect, err := cli.ContainerInspect(context.Background(), c.ID)
			if err != nil {
				logrus.Errorf("Error inspecting %s: %v\n", c.ID[:12], err)
				continue
			}

			var ips []string
			for _, net := range inspect.NetworkSettings.Networks {
				ips = append(ips, net.IPAddress)
			}
			numAddresses := len(ips)
			if numAddresses == 0 {
				logrus.Warnf("%s: No IP addresses found\n", c.Names[0])
				continue
			}

			logrus.Debugf("f %v ip\n", ips)

			exposedPorts := inspect.NetworkSettings.Ports
			if len(exposedPorts) == 0 {
				logrus.Warnf("%s: No exposed ports\n", c.Names[0])
				continue
			}

			logrus.Debugf("f %v exposedPorts\n", exposedPorts)

			var firstPort string
			for port := range exposedPorts {
				firstPort = port.Port()
				break
			}

			logrus.Debugf("f %v %v:%v send for worker\n", c.Names[0], ips, firstPort)

			go func(addresses []string) {
				for j := range addresses {
					jobs <- fmt.Sprintf("%s:%s", addresses[j], firstPort)
				}
			}(ips)

		}

	}
}
