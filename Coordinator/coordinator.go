/*
21I-0470 Aliza Ibrahim
21I-0603 Hamna Sadia Rizwan
21I-0572 Kissa Zahra //insection D
*/
package main

import (
	"crypto/tls"
	"fmt"
	"net/rpc"
	"sync"
)

type Matrix struct {
	Data [][]int
}

type MatrixOperation struct {
	A Matrix
	B Matrix
}

type Coordinator struct {
	mu      sync.Mutex
	workers []string
	next    int
}

func (c *Coordinator) RegisterWorker(workerAddr string, reply *string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.workers = append(c.workers, workerAddr)
	*reply = "Worker registered: " + workerAddr
	fmt.Println(*reply)
	return nil
}

func (c *Coordinator) Ping(args string, reply *string) error {
	*reply = "pong"
	return nil
}

func (c *Coordinator) assignWorker() (string, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if len(c.workers) == 0 {
		return "", fmt.Errorf("no workers available")
	}

	worker := c.workers[c.next]
	c.next = (c.next + 1) % len(c.workers)
	return worker, nil
}
func (c *Coordinator) connectToWorker(worker string) (*rpc.Client, error) {
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
	}

	conn, err := tls.Dial("tcp", worker, tlsConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to worker %s: %v", worker, err)
	}

	return rpc.NewClient(conn), nil
}

func (c *Coordinator) Add(args MatrixOperation, reply *Matrix) error {
	for retries := 0; retries < 3; retries++ {
		worker, err := c.assignWorker()
		if err != nil {
			return err
		}

		client, err := c.connectToWorker(worker)
		if err != nil {
			fmt.Println("Worker failed, reassigning...")
			continue
		}
		defer client.Close()

		return client.Call("MatrixService.Add", args, reply)
	}
	return fmt.Errorf("all workers failed for addition operation")
}

func (c *Coordinator) Transpose(args Matrix, reply *Matrix) error {
	for retries := 0; retries < 3; retries++ {
		worker, err := c.assignWorker()
		if err != nil {
			return err
		}

		client, err := c.connectToWorker(worker)
		if err != nil {
			fmt.Println("Worker failed, reassigning...")
			continue
		}
		defer client.Close()

		return client.Call("MatrixService.Transpose", args, reply)
	}
	return fmt.Errorf("all workers failed for transpose operation")
}

func (c *Coordinator) Multiply(args MatrixOperation, reply *Matrix) error {
	for retries := 0; retries < 3; retries++ {
		worker, err := c.assignWorker()
		if err != nil {
			return err
		}

		client, err := c.connectToWorker(worker)
		if err != nil {
			fmt.Println("Worker failed, reassigning...")
			continue
		}
		defer client.Close()

		return client.Call("MatrixService.Multiply", args, reply)
	}
	return fmt.Errorf("all workers failed for multiplication operation")
}

// establish connection bw worker node and coordinator
func startWorkerListener(coordinator *Coordinator) {
	cert, err := tls.LoadX509KeyPair("cert.pem", "key.pem") //load the tls cert & key
	if err != nil {
		fmt.Println("Error loading certificate:", err)
		return
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		MinVersion:   tls.VersionTLS12,
	}

	listener, err := tls.Listen("tcp", ":1234", tlsConfig)
	if err != nil {
		fmt.Println("Error starting TLS server:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Coordinator RPC TLS Server is running on port 1234 (for workers)...")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("TLS Connection error:", err)
			continue
		}
		go rpc.ServeConn(conn)
	}
}

// establish connection bw client and server
func startClientListener(coordinator *Coordinator) {
	cert, err := tls.LoadX509KeyPair("cert.pem", "key.pem")
	if err != nil {
		fmt.Println("Error loading certificate:", err)
		return
	}

	tlsConfig := &tls.Config{
		Certificates:       []tls.Certificate{cert},
		MinVersion:         tls.VersionTLS12,
		InsecureSkipVerify: true,
	}

	listener, err := tls.Listen("tcp", ":1235", tlsConfig) //port 1235
	if err != nil {
		fmt.Println("Error starting TLS server:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Coordinator RPC TLS Server is running on port 1235 (for clients)...")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("TLS Connection error:", err)
			continue
		}
		go rpc.ServeConn(conn)
	}
}
func main() {
	coordinator := new(Coordinator)
	rpc.Register(coordinator)

	go startWorkerListener(coordinator) //workers connc here
	go startClientListener(coordinator) //clients connc here

	select {}
}
