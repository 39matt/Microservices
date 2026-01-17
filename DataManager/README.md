# Data Manager
DataManager is a gRPC microservice written in Go that provides CRUD operations and aggregations (min, max, avg, sum) for IoT sensor readings stored in PostgreSQL.

## protoc command
Compiles the .proto file 
```shell
protoc --proto_path=internal/proto \
--go_out=internal/pb --go_opt=paths=source_relative \
--go-grpc_out=internal/pb --go-grpc_opt=paths=source_relative \
internal/proto/Reading.proto
```

## Docker instructions

### Postgres
1. docker run -d --name postgres --network iot-net -p 5433:5432 \
   -e POSTGRES_USER=user -e POSTGRES_PASSWORD=password -e POSTGRES_DB=iotdb \
   postgres:16
   - change user/pass for authentication on your local postgres
   - 5433 locally, 5432 in container
2. Run commands:
   - Create a database
    ```shell
    docker exec -it postgres psql -U matija -d postgres -c "CREATE DATABASE iotdb;"
    ```
   - Create a table inside it
    ```shell
    docker exec -it postgres psql -U matija -d iotdb -c "
    DROP TABLE IF EXISTS readings CASCADE;
    CREATE TABLE public.readings (
        id integer PRIMARY KEY,
        timestamp text NOT NULL,
        device_id text,
        co double precision,
        humidity real,
        light boolean,
        lpg double precision,
        motion boolean,
        smoke double precision,
        temperature real
    );"
    ```
   - Insert one 
    ```shell
    docker exec -it postgres psql -U matija -d iotdb -c "
    INSERT INTO readings (id, timestamp, device_id, co, humidity, light, lpg, motion, smoke, temperature)
    VALUES (1, '2026-01-16T17:27:00Z', 'device-001', 1.5, 60.2, true, 2.718, false, 0.3, 23.1)
    ON CONFLICT (id) DO NOTHING;"
    ```

### Server
1.  docker build -t datamanager .
2.  docker run -d --name datamanager --network iot-net -p 8081:8080 -e POSTGRES_URL="postgres://user:password@postgres:5432/iotdb?sslmode=disable" datamanager
    - connect to network (same as database)
    - expose 8080 on the container and bind it to 8081 locally
    - environment variable