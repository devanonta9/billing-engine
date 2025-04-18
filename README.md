# billing-engine
Repo for Billing Engine <br /><br />

# Requirements
Billing Engine (System Design and abstraction)
We are building a billing system for our Loan Engine. Basically the job of a billing engine is to provide the: <br />  
1. Loan schedule for a given loan(when am i supposed to pay how much) <br />
2. Outstanding Amount for a given loan <br />
3. Status of weather the customer is Delinquent or not <br /><br />

# Command
Run ```go mod tidy``` and ```go mod vendor``` to avoid inconsistency vendoring after clone the repo <br />
We can run it by using ```go run main.go``` or we can build the binary first <br />
To build binary run the command ```go build -o [name.exe] main.go``` if Windows or ```go build -o [name] main.go``` if macOS <br />
Then run the binary by using command ```./name.exe``` or ```./name``` <br /><br />

# Endpoint
Endpoint: <br />
1. Get User: Get ```baseUrl:portAddress/v1/users``` eg: ```localhost:8099/v1/users```. <br />
2. Create Loan: POST ```baseUrl:portAddress/v1/loans``` eg: ```localhost:8099/v1/loans```. POST consists data on the raw request body as json<br />
3. Check Delinquent: GET ```baseUrl:portAddress/v1/delinquents``` eg: ```localhost:8099/v1/delinquents```. GET consists data on the param<br />
4. Get Billing: GET ```baseUrl:portAddress/v1/billings``` eg: ```localhost:8099/v1/billings```. GET consists data on the param<br />
5. Create Payment: POST ```baseUrl:portAddress/v1/payments``` eg: ```localhost:8099/v1/payments```. POST consists data on the raw request body as json<br />
6. Get Schedule: GET ```baseUrl:portAddress/v1/schedules``` eg: ```localhost:8099/v1/schedules```. GET consists data on the param<br />

# Config
Config is read from config.yaml file (rename to config.yaml from config-example.yaml or create a new one) <br />
The default are: <br /><br />

```APP:``` <br />
 ```NAME: "app-name"``` <br />
 ```PORT: "app-port"``` <br />
```ROUTE:``` <br />
  ```METHODS: "[GET, POST, PUT, DELETE]"``` <br />
  ```HEADERS: "[Content-Type]"``` <br />
  ```ORIGIN: "['*']"``` <br />
```DATABASE:``` <br />
  ```READ:``` <br />
    ``USERNAME: "db-username"`` <br />
    ``PASSWORD: "db-password"`` <br />
    ``URL: "db-host"`` <br />
    ``PORT: "db-port"`` <br />
    ``DB_NAME: "db-name"`` <br />
    ``MAXIDLECONNS: "5"`` <br />
    ``MAXOPENCONNS: "5"`` <br />
    ``MAXLIFETIME: "31"`` <br />
    ``TIMEOUT: "10"`` <br />
    ``SSL_MODE: "disable"`` <br /><br />

# Collection
Postman and Script are available in the repository <br /><br />

# Tech Stack 
1. Golang <br />
2. PostgreSQL