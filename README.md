# TestON
TestON is a Discord bot that provides users with recent changes and statistics from Wikipedia. The bot is designed to store raw Wikipedia events in a PostgreSQL database, retrieve recent changes, and display statistics based on language and date. It is built in Go. 

## Features
- **Set Language Preference:**  
  Users can set their default language for Wikipedia articles.  
- **Fetch Recent Changes:**  
  Retrieve a specified number of recent Wikipedia edits in the user’s preferred language (with a configurable limit, up to 100).  
- **Containerized Deployment:**  
  Includes Docker Compose configuration for streamlined local development and deployment.
- **View Statistics:**  
  Display the number of changes for a particular language on a given date.  

## Prerequisites
- **Go:**
  Version 1.20 or later    
- **PostgreSQL:**  
  A PostgreSQL database instance. Make sure there is nothing running locally on your machine on port 5432:5432.  
- **Docker:**  
  For running PostgreSQL
- **Make:**  
  To run the migrations
- **Discord Bot Token:**  
  Obtain token by creating a bot in the Discord Developer Portal

## Installation

- **Go:**
  Version 1.20 or later    
- **PostgreSQL:**  
  A PostgreSQL database instance. Make sure there is nothing running locally on your machine on port 5432:5432.  
- **Docker:**  
  For running PostgreSQL  
- **Discord Bot Token:**  
  Obtain token by creating a bot in the Discord Developer Portal

## Running the Program

1. **Clone the Repository:**

   ```bash
   git clone https://github.com/vlkhvnn/TestON.git
   cd TestON
2. **Create .env file in the root directory with the following variables**  

   ```bash
   DISCORD_TOKEN=your_discord_bot_token
   DB_ADDR=postgres://postgres:1234@localhost:5432/teston?sslmode=disable
   DB_MAX_OPEN_CONNS=30
   DB_MAX_IDLE_CONNS=30
   DB_MAX_IDLE_TIME=15m
   ```
   Make sure to get a bot token from discord developers site and paste it in DISCORD_TOKEN field.

3. **Running the Application:**

   Ensure Docker is running on your system, then execute:
   ```bash
   docker-compose up
   ```
   On the another terminal run:
   ```bash
   cd cmd/
   go run .
   ```
4. **Run the Migrations**
   On the another terminal run. (If you do not have golang-migrate):
   ```bash
   brew install golang-migrate
   ```
   ```bash
   make migrate-up
   ```

## Usage  
Now your bot is ready to work. Invite your bot to the server.
- **Set Language Preference:**
  ```bash
  !setLang [language_code]
  !setLang en
  !setLang es
  !setLang ru
  ```
- **Fetch Recent Changes:**
  ```bash
  !recent [optional: number_of_events]
  !recent 5
  !recent en 20
  ```
- **View Statistics:**
  ```bash
  !stats [yyyy-mm-dd] [optional: language_code]
  !stats 2025-02-04 en
  ```

## Scaling Architecture for Higher Throughput

For higher volumes of Wikipedia events, consider integrating additional technologies:

1. **Apache Kafka for Event Ingestion:**
   - **Publish Raw Wikipedia Events:**  
     Instead of processing events directly in the bot, publish raw events to a Kafka topic.
   - **High Throughput:**  
     Kafka can handle massive volumes of data and provides durability and scalability.

2. **Apache Spark for Stream Processing:**
   - **Process Events in Real-Time:**  
     Use Spark Streaming to consume events from Kafka. Spark can filter, aggregate, and transform data in real-time.
   - **Output to Database or Another Kafka Topic:**  
     After processing, Spark can write the filtered/aggregated events back to a database (e.g., PostgreSQL) or another Kafka topic.

3. **Bot Consumption:**
   - **Consume Processed Data:**  
     The Discord bot can consume data from the database (or directly from a Kafka topic if needed) for more efficient responses.

4. **Horizontal Scaling:**
   - **Scale Kafka:**  
     Increase the number of brokers and partitions for parallel processing.
   - **Scale Spark:**  
     Add more worker nodes to distribute processing load.
   - **Scale Bot Instances:**  
     Deploy multiple instances of the bot behind a load balancer if necessary.

5. **Monitoring and Auto-scaling:**
   - **Use Monitoring Tools:**  
     Tools like Prometheus and Grafana can monitor system performance.
   - **Auto-scaling:**  
     Use container orchestration (e.g., Kubernetes) to auto-scale services based on load.

   
