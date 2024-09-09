# August 6 Notes

## Explain Docker-compose for kafka
Here's a breakdown of your Docker Compose file:

### Version
```yaml
version: "3"
```
- Specifies the version of the Docker Compose file format.

### Services
- Defines a list of services (containers) to be run.

#### Kafka Service
```yaml
services:
  kafka:
    image: 'bitnami/kafka:latest'
    ports:
      - "9092:9092"
    volumes:
      - "kafka_data:/bitnami"
    environment:
```
- **image:** Specifies the Docker image to use for Kafka (`bitnami/kafka:latest`).
- **ports:** Maps port 9092 on the host to port 9092 on the Kafka container.
- **volumes:** Mounts a named volume (`kafka_data`) to `/bitnami` inside the container for persistent storage.

#### Environment Variables
These variables configure Kafka settings:

```yaml
      # KRaft settings
      - KAFKA_CFG_NODE_ID=0
      - KAFKA_CFG_PROCESS_ROLES=controller,broker
      - KAFKA_CFG_CONTROLLER_QUORUM_VOTERS=0@kafka:9093
```
- **KAFKA_CFG_NODE_ID:** Unique ID for the Kafka node.
- **KAFKA_CFG_PROCESS_ROLES:** Specifies the roles of the node (both controller and broker).
- **KAFKA_CFG_CONTROLLER_QUORUM_VOTERS:** Configures the controller quorum, with `0@kafka:9093` indicating node 0 is the controller and will listen on port 9093.

#### Listeners Configuration
```yaml
      # Listeners
      - KAFKA_CFG_LISTENERS=PLAINTEXT://:9092,CONTROLLER://:9093
      - KAFKA_CFG_ADVERTISED_LISTENERS=PLAINTEXT://localhost:9092
      - KAFKA_CFG_LISTENER_SECURITY_PROTOCOL_MAP=CONTROLLER:PLAINTEXT,PLAINTEXT:PLAINTEXT
      - KAFKA_CFG_CONTROLLER_LISTENER_NAMES=CONTROLLER
      - KAFKA_CFG_INTER_BROKER_LISTENER_NAME=PLAINTEXT
```
- **KAFKA_CFG_LISTENERS:** Defines listeners for different protocols and ports:
  - `PLAINTEXT://:9092` for standard Kafka communication.
  - `CONTROLLER://:9093` for controller communication.
- **KAFKA_CFG_ADVERTISED_LISTENERS:** Defines how this Kafka broker will advertise itself to clients (in this case, as `localhost:9092`).
- **KAFKA_CFG_LISTENER_SECURITY_PROTOCOL_MAP:** Maps listener names to security protocols (both are `PLAINTEXT`).
- **KAFKA_CFG_CONTROLLER_LISTENER_NAMES:** Specifies the listener names for the controller role (`CONTROLLER`).
- **KAFKA_CFG_INTER_BROKER_LISTENER_NAME:** Defines the listener name for inter-broker communication (`PLAINTEXT`).

### Volumes
```yaml
volumes:
  kafka_data:
    driver: local
```
- Defines a named volume (`kafka_data`) with the local driver for persistent storage of Kafka data.

### Summary
This Docker Compose file sets up a Kafka container using the Bitnami Kafka image, configures it to act as both a broker and controller using the KRaft mode, and sets up necessary listeners and advertised listeners for communication. The `kafka_data` volume ensures that Kafka's data persists across container restarts.


## Explanation of "PLAINTEXT"
In the context of your Docker Compose configuration for Kafka, "PLAINTEXT" refers to the communication protocol that Kafka uses for its network traffic. Here's a more detailed explanation of where "PLAINTEXT" is used in your configuration:

### Contexts of "PLAINTEXT"

1. **KAFKA_CFG_LISTENERS**
   ```yaml
   - KAFKA_CFG_LISTENERS=PLAINTEXT://:9092,CONTROLLER://:9093
   ```
   - **PLAINTEXT://:9092:** This specifies that Kafka will listen for client connections on port 9092 using the PLAINTEXT protocol, which means the data is transmitted without any encryption.
   - **CONTROLLER://:9093:** This specifies that Kafka's controller will listen on port 9093, but it's mapped to the PLAINTEXT protocol later in the configuration.

2. **KAFKA_CFG_ADVERTISED_LISTENERS**
   ```yaml
   - KAFKA_CFG_ADVERTISED_LISTENERS=PLAINTEXT://localhost:9092
   ```
   - **PLAINTEXT://localhost:9092:** This tells Kafka to advertise itself to clients using the PLAINTEXT protocol on `localhost:9092`.

3. **KAFKA_CFG_LISTENER_SECURITY_PROTOCOL_MAP**
   ```yaml
   - KAFKA_CFG_LISTENER_SECURITY_PROTOCOL_MAP=CONTROLLER:PLAINTEXT,PLAINTEXT:PLAINTEXT
   ```
   - This maps the `CONTROLLER` listener to use the PLAINTEXT protocol and confirms that the `PLAINTEXT` listener will use the PLAINTEXT protocol.

4. **KAFKA_CFG_INTER_BROKER_LISTENER_NAME**
   ```yaml
   - KAFKA_CFG_INTER_BROKER_LISTENER_NAME=PLAINTEXT
   ```
   - This specifies that the inter-broker communication will use the PLAINTEXT listener, meaning that the communication between different Kafka brokers will also be unencrypted.

### Summary
- **PLAINTEXT Protocol:** This is the simplest form of communication protocol where data is transmitted in plain text without any encryption. It is easy to set up and use, but it does not provide any security features such as encryption or authentication.

In a production environment, using PLAINTEXT means that anyone who can intercept the network traffic can read the data. Therefore, for secure deployments, you might want to consider using SSL/TLS encryption (`SSL` protocol) to ensure the confidentiality and integrity of the data being transmitted.

## Difference between controller, broker and PLAINTEXT

In the context of Apache Kafka and your Docker Compose configuration, "controller," "broker," and "PLAINTEXT" have specific meanings related to the roles of Kafka nodes and the communication protocols used. Here's a detailed explanation:

### Controller

- **Role:** The controller in Kafka is responsible for managing the cluster metadata, including partition leadership and replica assignments. It ensures that the cluster operates correctly and coordinates changes in the cluster topology, such as adding or removing brokers and reassigning partitions.
- **Responsibilities:**
  - Elects partition leaders.
  - Monitors the health of brokers.
  - Manages topic and partition metadata.
- **Configuration in Docker Compose:**
  ```yaml
  - KAFKA_CFG_PROCESS_ROLES=controller,broker
  ```
  This indicates that the Kafka node will act as both a controller and a broker.

### Broker

- **Role:** A Kafka broker is a server that stores and serves Kafka topics. Brokers handle client requests for reading and writing data, as well as managing the replication of data across the cluster.
- **Responsibilities:**
  - Stores data (topics) on disk.
  - Serves client read and write requests.
  - Replicates data to other brokers for fault tolerance.
- **Configuration in Docker Compose:**
  ```yaml
  - KAFKA_CFG_PROCESS_ROLES=controller,broker
  ```
  This indicates that the Kafka node will act as both a broker and a controller.

### PLAINTEXT

- **Definition:** PLAINTEXT is a communication protocol used by Kafka for transmitting data without encryption. It means that data sent between clients and brokers or between brokers is not encrypted and can be intercepted and read by anyone with access to the network.
- **Usage in Configuration:**
  - **Listeners:**
    ```yaml
    - KAFKA_CFG_LISTENERS=PLAINTEXT://:9092,CONTROLLER://:9093
    ```
    This specifies that Kafka will listen for client connections on port 9092 using the PLAINTEXT protocol and for controller communications on port 9093.
  - **Advertised Listeners:**
    ```yaml
    - KAFKA_CFG_ADVERTISED_LISTENERS=PLAINTEXT://localhost:9092
    ```
    This tells Kafka to advertise itself to clients using the PLAINTEXT protocol on `localhost:9092`.
  - **Listener Security Protocol Map:**
    ```yaml
    - KAFKA_CFG_LISTENER_SECURITY_PROTOCOL_MAP=CONTROLLER:PLAINTEXT,PLAINTEXT:PLAINTEXT
    ```
    This maps both the `CONTROLLER` and `PLAINTEXT` listeners to use the PLAINTEXT protocol.
  - **Inter-Broker Listener:**
    ```yaml
    - KAFKA_CFG_INTER_BROKER_LISTENER_NAME=PLAINTEXT
    ```
    This specifies that inter-broker communication will use the PLAINTEXT listener.

### Summary
- **Controller:** Manages cluster metadata and coordinates cluster operations.
- **Broker:** Stores data, handles client requests, and replicates data.
- **PLAINTEXT:** A protocol for unencrypted data transmission, easy to set up but lacks security features like encryption and authentication.

## Ports of 9092 and 9093

In your Docker Compose configuration for Kafka, ports 9092 and 9093 serve distinct purposes related to the different roles and communication needs within a Kafka cluster. Here's an explanation of what each port is used for:

### Port 9092

- **Purpose:** Client Communication
- **Usage:** Port 9092 is used for communication between Kafka clients (such as producers and consumers) and the Kafka broker. When clients want to produce messages to or consume messages from Kafka topics, they connect to this port.
- **Configuration:**
  ```yaml
  - KAFKA_CFG_LISTENERS=PLAINTEXT://:9092
  - KAFKA_CFG_ADVERTISED_LISTENERS=PLAINTEXT://localhost:9092
  ```
  - `KAFKA_CFG_LISTENERS=PLAINTEXT://:9092`: This configuration tells Kafka to listen for client connections on port 9092 using the PLAINTEXT protocol.
  - `KAFKA_CFG_ADVERTISED_LISTENERS=PLAINTEXT://localhost:9092`: This configuration advertises the broker's address to clients, indicating that they should connect to `localhost:9092` to interact with the Kafka broker.

### Port 9093

- **Purpose:** Controller Communication
- **Usage:** Port 9093 is used for communication between Kafka controllers and other brokers. This port is dedicated to controller-related traffic, which includes the management of cluster metadata and coordination tasks. In a KRaft (Kafka Raft) mode setup, the controller role is especially significant as it ensures proper cluster operations.
- **Configuration:**
  ```yaml
  - KAFKA_CFG_LISTENERS=CONTROLLER://:9093
  - KAFKA_CFG_CONTROLLER_QUORUM_VOTERS=0@kafka:9093
  - KAFKA_CFG_LISTENER_SECURITY_PROTOCOL_MAP=CONTROLLER:PLAINTEXT
  - KAFKA_CFG_CONTROLLER_LISTENER_NAMES=CONTROLLER
  ```
  - `KAFKA_CFG_LISTENERS=CONTROLLER://:9093`: This configuration tells Kafka to listen for controller communication on port 9093.
  - `KAFKA_CFG_CONTROLLER_QUORUM_VOTERS=0@kafka:9093`: This specifies that the node with ID 0 (your Kafka node) is a controller and will communicate on port 9093.
  - `KAFKA_CFG_LISTENER_SECURITY_PROTOCOL_MAP=CONTROLLER:PLAINTEXT`: This maps the controller listener to use the PLAINTEXT protocol.
  - `KAFKA_CFG_CONTROLLER_LISTENER_NAMES=CONTROLLER`: This specifies that the listener name for the controller role is `CONTROLLER`.

### Summary

- **Port 9092:** Used for client-broker communication (e.g., producers and consumers connecting to the Kafka broker). Clients connect to this port to send and receive messages from Kafka.
- **Port 9093:** Used for controller-broker communication (e.g., managing cluster metadata and coordinating tasks between brokers). This port is specifically for internal Kafka operations, especially in a KRaft mode setup where the controller plays a crucial role.

## Cluster and Cluster Metadata

### Kafka Cluster and Cluster Metadata

#### Kafka Cluster
- **Definition:** A Kafka cluster is a group of Kafka broker nodes working together to handle data streaming and messaging. The cluster provides redundancy, scalability, and fault tolerance.
- **Components:**
  - **Brokers:** These are the nodes (servers) in the Kafka cluster. Each broker stores data and serves client requests.
  - **Topics:** Logical channels to which data is written and from which data is read. Each topic can have multiple partitions.
  - **Partitions:** A topic is divided into partitions to allow parallel processing. Each partition is an ordered, immutable sequence of records.
  - **Producers:** Clients that send (produce) data to Kafka topics.
  - **Consumers:** Clients that read (consume) data from Kafka topics.

#### Cluster Metadata
- **Definition:** Cluster metadata refers to the information that Kafka maintains about the structure and state of the cluster. This includes details about:
  - **Brokers:** IDs and addresses of all brokers in the cluster.
  - **Topics:** Names of all topics in the cluster.
  - **Partitions:** Number of partitions for each topic and their assignment to brokers.
  - **Leaders and Replicas:** Information about which broker is the leader for each partition and the replicas of that partition.
  - **Offsets:** The position of the last consumed record in each partition for each consumer group.

- **Importance:** Cluster metadata is crucial for the operation of a Kafka cluster. It allows producers to send data to the correct broker and partition, consumers to read data from the correct broker and partition, and brokers to coordinate replication and failover.

### Kafka Nodes

- **Definition:** A Kafka node, also known as a broker, is a single server within a Kafka cluster. Each node can serve multiple roles (broker, controller) and works together with other nodes to provide the cluster's functionality.
- **Roles:**
  - **Broker:** Stores data (topics) and serves client requests for data reads and writes.
  - **Controller:** Manages the metadata for the cluster, including partition leadership and replication. The controller coordinates with other brokers to ensure data is properly replicated and that the cluster operates smoothly.
  - **In KRaft Mode:** In the KRaft (Kafka Raft) mode, the controller role is integrated more tightly with the broker, and the coordination tasks are managed through a Raft consensus algorithm.

### Summary

- **Cluster:** A group of Kafka broker nodes working together to handle data streaming and messaging tasks.
- **Cluster Metadata:** Information about the brokers, topics, partitions, leaders, replicas, and offsets in the cluster. This metadata is crucial for the correct operation and coordination of the cluster.
- **Kafka Nodes (Brokers):** Individual servers in the Kafka cluster that store data and serve client requests. They also participate in managing the cluster's state and ensuring data replication and failover.