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

func printMatrix(matrix Matrix) {
	for _, row := range matrix.Data {
		fmt.Println(row)
	}
	fmt.Println()
}

func main() {
	tlsConfig := &tls.Config{InsecureSkipVerify: true}

	conn, err := tls.Dial("tcp", "0.tcp.in.ngrok.io:19142", tlsConfig) //this is connected to port 1235
	if err != nil {
		fmt.Println("Error connecting to coordinator:", err)
		return
	}
	defer conn.Close()

	client := rpc.NewClient(conn)
	defer client.Close()

	matrixA := Matrix{
		Data: [][]int{
			{1, 2, 3},
			{4, 5, 6},
			{7, 8, 9},
		},
	}

	matrixB := Matrix{
		Data: [][]int{
			{14, 8, 7},
			{6, 5, 4},
			{3, 2, 1},
		},
	}

	//matrix operationss
	addArgs := MatrixOperation{A: matrixA, B: matrixB} //addition
	var addResult Matrix
	err = client.Call("Coordinator.Add", addArgs, &addResult)
	if err != nil {
		fmt.Println("RPC Addition error:", err)
	} else {
		fmt.Println("Matrix Addition Result:")
		printMatrix(addResult)
	}

	var transposeResult Matrix //transpose
	err = client.Call("Coordinator.Transpose", matrixA, &transposeResult)
	if err != nil {
		fmt.Println("RPC Transpose error:", err)
	} else {
		fmt.Println("Matrix Transpose Result (A):")
		printMatrix(transposeResult)
	}

	mulArgs := MatrixOperation{A: matrixA, B: matrixB} //multiplication
	var mulResult Matrix
	err = client.Call("Coordinator.Multiply", mulArgs, &mulResult)
	if err != nil {
		fmt.Println("RPC Multiplication error:", err)
	} else {
		fmt.Println("Matrix Multiplication Result:")
		printMatrix(mulResult)
	}
}
