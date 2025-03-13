# 🔢 RPC Matrix Server  

To check the assignment question, click here: **[CS4049_Assignment_01](CS4049_Assignment_01)** 

## 🔐 Generate TLS Certificates  
Run:  
```sh
openssl req -x509 -newkey rsa:2048 -keyout key.pem -out cert.pem -days 365 -nodes
```

##  Put cert.pem and key.pem in:
📂 client/      
📂 coordinator/  
📂 worker/      


##  1️⃣ Start the coordinator (server):

```sh
python coordinator/coordinator.py
```
##  2️⃣ Start workers (each in a separate terminal):

```sh
python worker/worker.py
```
## 3️⃣ Run the client:
```sh
python client/client.py
```
