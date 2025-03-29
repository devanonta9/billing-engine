# billing-engine
Repo for Billing Engine <br /><br />

# Requirements
Billing Engine (System Design and abstraction)
We are building a billing system for our Loan Engine. Basically the job of a billing engine is to provide the: <br />  
1. Loan schedule for a given loan(when am i supposed to pay how much) <br />
2. Outstanding Amount for a given loan <br />
3. Status of weather the customer is Delinquent or not <br /><br />

# Command
We can run it by using ```go run main.go``` or we can build the binary first <br />
To build binary run the command ```go build -o [name.exe] main.go``` if Windows or ```go build -o [name] main.go``` if macOS <br />
Then run the binary by using command ```./name.exe``` or ```./name``` <br />

# Endpoint
Auth Endpoint: <br />
1. Login: POST ```baseUrl:portAddress/auth/login``` eg: ```localhost:8099/auth/login```. POST consists 2 data ```username``` and ```password``` on the raw request body as json <br />
2. Logout: GET ```baseUrl:portAddress/auth/logout``` eg: ```localhost:8099/auth/logout```. <br /><br />

Store Endpoint: <br />
1. Upload File: POST ```baseUrl:portAddress/store/upload``` eg: ```localhost:8099/store/upload```. POST consist 1 data ```file``` on the form request body as file <br />
2. Get List Files: GET ```baseUrl:portAddress/store/list``` eg: ```localhost:8099/store/list```. <br />
3. Download File: GET ```baseUrl:portAddress/store/download``` eg: ```localhost:8099/store/download```. <br /><br />

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
    ``SSL_MODE: "disable"`` <br />

# Collection
Postman API collection = ```https://www.getpostman.com/collections/82b618929a79a13df4ce``` <br />