# Internet of Things and Services – Project

This repository contains multiple projects built as part of an **Internet of Things and Services** course. The system demonstrates a microservice-based architecture for managing, processing, and visualizing environmental sensor data.

Each project is containerized using **Docker Compose** and connected through a common network for seamless interoperability.

## 🧱 Overall Architecture
```mermaid
---
config:
  theme: base
  themeVariables:
    primaryColor: '#7cd4fc'
    primaryTextColor: '#000'
    primaryBorderColor: '#333'
    lineColor: '#666'
    fontSize: 14px
  layout: elk
---
flowchart LR
    CSV["iot_telemetry_data.csv"] --> SG["SensorGenerator<br>Python Script"]
    SG -- REST POST /readings --> GW["Gateway<br>.NET REST API<br>:5236"]
    GW -- gRPC --> DM["DataManager<br>Go gRPC Server<br>:8080"]
    DM -- SQL --> DB[("Postgres<br>iotdb.readings<br>:5432")]
    DM -. PUBLISH /reading .-> MQTT["MQTT Broker<br>Mosquitto<br>:1883"]
    MQTT -. SUBSCRIBE /reading .-> EM["Event Manager<br>Go MQTT Handler<br>:8090"] & AL["Analytics"]
    EM -. PUBLISH /limit .-> MQTT
    AL -- REST --> ML["MLaaS<br>Python FastAPI<br>:8000"]
    AL -. PUBLISH /example .-> NATS["NATS Broker<br>Example<br>:1884"]
    MQTT -. SUBSCRIBE /limit .-> CLIENT["MQTT &amp; NATS Client<br>NextJS<br>:3000"]
    NATS -. SUBSCRIBE /example .-> CLIENT
    PM["Postman/grpcurl"] -. HTTPS .-> GW
    PM -. gRPC .-> DM

     SG:::generator
     GW:::service
     DM:::service
     DB:::db
     EM:::service
     MQTT:::broker
     NATS:::broker
     AL:::service
     ML:::service
     CLIENT:::service
     PM:::test
    classDef service fill:#7cd4fc,stroke:#333,stroke-width:2px,color:#000
    classDef db fill:#a3f58e,stroke:#333,stroke-width:2px,color:#000
    classDef broker fill:#ba52fa,stroke:#333,stroke-width:2px,color:#fff
    classDef generator fill:#fadb52,stroke:#333,stroke-width:2px,color:#000
    classDef test stroke:#ff9800,stroke-width:3px,stroke-dasharray: 5 5
```
The entire system is structured as a set of loosely coupled microservices, each responsible for a specific task:

- **Data Manager:** gRPC-based data service written in Go.  
- **Gateway:** ASP.NET Core MVC application that acts as an API gateway.  
- **Sensor Generator:** Python script that simulates IoT sensor data.  

A structural diagram of the system (coming soon) illustrates how containers communicate using both HTTP and gRPC protocols.

The dataset used for simulation is sourced from [Environmental Sensor Data on Kaggle](https://www.kaggle.com/datasets/garystafford/environmental-sensor-data-132k).


## 📦 Project 1 – Environmental Data System (Gateway, DataManager and Database)
### 1. Data Manager (Go)

A **gRPC service** responsible for data storage and CRUD operations.

- Connects to **PostgreSQL**.
- Exposes methods over **gRPC** (runs on HTTP/2).
- Listens on a TCP port to handle remote procedure calls from other services.
- Implements protobuf definitions for data models and RPC methods.

### 2. Gateway (ASP.NET Core 10 MVC)

An **API Gateway and frontend** that communicates with the Data Manager via gRPC.

- Uses **OpenAPI** for [specification.](https://app.swaggerhub.com/apis/elfak-695/Project1/1.0.0)  
- Implements **MVC architecture** for a clean separation of concerns.  
- Provides a REST interface for clients to access sensor data.  
- Forwards data operations to the Data Manager via gRPC calls.

### 3. Sensor Generator (Python)

A lightweight **data simulator** that mimics IoT devices by sending sensor readings.

- Reads environmental data from the Kaggle CSV file.  
- Randomly selects rows and sends HTTP requests to the Gateway API.  
- Designed for testing and system load simulation.

---

## 🐳 Docker Setup

All services are containerized and orchestrated using **Docker Compose**.

- Each component (Gateway, Data Manager, PostgreSQL) runs in its own container.
- Shared network **iot-net** enables cross-service communication.
- Environment variables handle ports, credentials, and connection strings.

## ⚙️ How to Run Locally

The simplest way to start all services together:

```bash
docker-compose up --build
```

This command will:
- Build and start the **PostgreSQL**, **Data Manager**, **Gateway**, and containers.
- Create a shared Docker network for service communication.
- Expose ports as configured inside `docker-compose.yml`.

Once started:
- Gateway should be available at: [**http://localhost:5237**](http://localhost:5237)
- Data Manager runs internally and communicates via gRPC.

To stop containers:

```bash
docker-compose down
```

---

## 📈 Future Projects

### Project 2 – MQTT
Addition of MQTT protocol using Mosquitto.
[AsyncAPI Specification](https://studio.asyncapi.com/?url=https://raw.githubusercontent.com/39matt/Microservices/refs/heads/main/DataManager/asyncapi.txt)

### Project 3 – TBD
_(Potential service for notification, ML predictions, or sensor anomaly detection)_

---

## 📚 Technologies Used

| Component | Language / Framework | Key Features |
|------------|---------------------|---------------|
| Data Manager | Go | gRPC, PostgreSQL, CRUD over RPC |
| Gateway | ASP.NET Core 10 | MVC, REST + gRPC integration, OpenAPI |
| Sensor Generator | Python | Data simulation, HTTP client |
| Infrastructure | Docker, Docker Compose | Containerization, orchestration |
