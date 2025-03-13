# 🔢 RPC Matrix Server  

To check the assignment question, click here: **[CS4049_Assignment_01](CS4049_Assignment_01.pdf)** 

## 🔐 Generate TLS Certificates  
Run:  
```sh
openssl req -x509 -newkey rsa:2048 -keyout key.pem -out cert.pem -days 365 -nodes
```

##  Put cert.pem and key.pem in:
📂 client/      
📂 coordinator/  
📂 worker/      


##  1️⃣ Start the coordinator:

```sh
go run Coordinator/coordinator.go
```
##  2️⃣ Start workers (each in a separate terminal):

```sh
go run worker/worker.go
```
## 3️⃣ Run the client:
```sh
go run client/client.go
```
