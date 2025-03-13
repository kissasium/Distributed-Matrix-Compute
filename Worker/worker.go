/*
21I-0470 Aliza Ibrahim
21I-0603 Hamna Sadia Rizwan
21I-0572 Kissa Zahra //in section D
*/
package main

import (
	"crypto/tls"

	"fmt"
	"net/rpc"
)

type Matrix struct {
	Data [][]int
}

type MatrixOperation struct {
	A Matrix
	B Matrix
}

type MatrixService struct{}

func (m *MatrixService) Add(args MatrixOperation, reply *Matrix) error {
	fmt.Println("Processing Addition on Worker")
	rows := len(args.A.Data)
	cols := len(args.A.Data[0])

	result := Matrix{Data: make([][]int, rows)}
	for i := range result.Data {
		result.Data[i] = make([]int, cols)
		for j := range result.Data[i] {
			result.Data[i][j] = args.A.Data[i][j] + args.B.Data[i][j]
		}
	}
	*reply = result
	fmt.Println("Addition Completed")
	return nil
}

func (m *MatrixService) Transpose(args Matrix, reply *Matrix) error {
	fmt.Println("Processing Transpose on Worker")
	rows := len(args.Data)
	cols := len(args.Data[0])

	result := make([][]int, cols)
	for i := range result {
		result[i] = make([]int, rows)
		for j := 0; j < rows; j++ {
			result[i][j] = args.Data[j][i]
		}
	}
	reply.Data = result
	fmt.Println("Transpose Completed")
	return nil
}

func (m *MatrixService) Multiply(args MatrixOperation, reply *Matrix) error {
	fmt.Println("Processing Multiplication on Worker")
	rowsA := len(args.A.Data)
	colsA := len(args.A.Data[0])
	rowsB := len(args.B.Data)
	colsB := len(args.B.Data[0])

	if colsA != rowsB {
		return fmt.Errorf("invalid matrix dimensions for multiplication")
	}

	result := make([][]int, rowsA)
	for i := range result {
		result[i] = make([]int, colsB)
		for j := 0; j < colsB; j++ {
			for k := 0; k < colsA; k++ {
				result[i][j] += args.A.Data[i][k] * args.B.Data[k][j]
			}
		}
	}
	reply.Data = result
	fmt.Println("Multiplication Completed")
	return nil
}

func main() {
	worker := new(MatrixService)
	rpc.Register(worker)

	cert, err := tls.LoadX509KeyPair("cert.pem", "key.pem") //load tls cert & key
	if err != nil {
		fmt.Println("Error loading certificate:", err)
		return
	}

	tlsConfig := &tls.Config{ //for worker listener
		Certificates:       []tls.Certificate{cert},
		MinVersion:         tls.VersionTLS12,
		InsecureSkipVerify: true,
	}

	listener, err := tls.Listen("tcp", ":0", tlsConfig) //for tls
	if err != nil {
		fmt.Println("Error starting worker:", err)
		return
	}
	defer listener.Close()

	workerAddr := listener.Addr().String()
	fmt.Println("Worker started on", workerAddr)

	conn, err := tls.Dial("tcp", ":1234", tlsConfig) //connc to coordinator using tls
	if err != nil {
		fmt.Println("Failed to connect to coordinator:", err)
		return
	}
	defer conn.Close()

	client := rpc.NewClient(conn)
	if client == nil {
		fmt.Println("Failed to create RPC client")
		return
	}
	defer client.Close()

	var pingReply string //test connc
	err = client.Call("Coordinator.Ping", "test", &pingReply)
	if err != nil {
		fmt.Println("Connection test failed:", err)
		return
	}

	var reply string
	err = client.Call("Coordinator.RegisterWorker", workerAddr, &reply)
	if err != nil {
		fmt.Println("Registration failed:", err)
		return
	}
	fmt.Println(reply)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Connection error:", err)
			continue
		}
		go rpc.ServeConn(conn)
	}
}
