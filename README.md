# 🚀 TiKV Admin Web UI

![Stream](https://img.shields.io/endpoint?url=https://gist.githubusercontent.com/HayesGordon/e7f3c4587859c17f3e593fd3ff5b13f4/raw/11d9d9385c9f34374ede25f6471dc743b977a914/badge.json)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/GetStream/tikv-ui)
![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)
![GitHub contributors](https://img.shields.io/github/contributors/GetStream/tikv-ui)
[![GitHub Sponsors](https://img.shields.io/github/sponsors/GetStream)](https://github.com/sponsors/GetStream)

## 🌟 Overview

**TiKV Admin Web UI is a full-stack web application and REST API for the powerful, intuitive management of TiKV clusters.**

Tired of juggling command-line tools? This project provides a dedicated, graphical frontend interface for exploring your data, managing cluster connections, and performing raw key-value manipulation with ease.

## ✨ Key Features

This application is built to be the essential administrative tool for your TiKV deployment:

- **🌐 Multi-Cluster Management:** Dynamically connect to and switch between multiple TiKV clusters without restarting the server. Perfect for managing development, staging, and production environments.
- **🔎 Data Explorer:** Perform powerful **Raw KV Operations** directly through the UI: Get, Put, Delete, and Range Scans.
- **🧠 Smart Value Parsing:** Automatically detects and parses complex value formats, including **MsgPack** and **JSON**, providing you with structured, readable data alongside the raw string value.
- **💻 Dedicated Web Frontend:** A modern interface built with TypeScript, offering a smooth and efficient way to interact with your key-value store.
- **🔌 Comprehensive REST API:** An easy-to-use API backend allows for integration into custom monitoring or automation workflows.

--

## 🚀 Getting Started

The TiKV Admin Web UI is built with Go for the backend API and a modern web stack for the frontend.

The easiest way to run the TiKV Admin Web UI is by using the pre-built Docker image.

### Option 1: Docker Compose (Recommended)

Use the official image in your `docker-compose.yml` file. This example assumes you have TiKV PD nodes named `pd0`, `pd1`, etc., running in the same network.

```yaml
version: "3.7"
services:
  # ... your TiKV cluster services (pd0, tikv0, etc.) ...

  tikv-ui:
    image: ghcr.io/getstream/tikv-ui:latest
    container_name: tikv-admin-ui
    environment:
      # **IMPORTANT:** Set TIKV_PD_ADDRS to the addresses of your PD nodes
      TIKV_PD_ADDRS: "pd0:2379,pd1:2379"
    # Use depends_on to ensure PD nodes are up before the UI starts
    depends_on:
      - "pd0"
      - "pd1"
      - "pd2" # Adjust this list to match your cluster
    ports:
      - "8081:8081"
```

### Option 2: Run Directly from Source

If you prefer to run the application directly, you will need Go installed.

### Prerequisites

- Go (v1.18+)
- A running TiKV cluster with access to its PD addresses.

### 🔨 Building the Application

Clone the repository and build the Go executable:

```bash
git clone https://github.com/GetStream/tikv-ui.git
cd tikv-ui
go build -o bin/tikv-ui ./cmd/tikv-ui
```

## 🏃 Running the app

Set the `TIKV_PD_ADDRS` environment variable with comma-separated PD addresses for your default cluster, then run the executable.
the format for `TIKV_PD_ADDRS` is {host}{,other hosts}{|cluster name}{;other clusters}

e.g.

```bash
export TIKV_PD_ADDRS="127.0.0.1:2379|default;pd-us-west-1.aws.com:2379|My US WEST cluster"
```

at least one host is required, other params are optional

```bash
# if you want to run on a specific port, just export the port variable, e.g. export PORT=8082
# Set your TiKV PD addresses (required)
export TIKV_PD_ADDRS="127.0.0.1:2379"

# Run the server
./bin/tikv-ui

# or run development server
make dev

# or separate backend
make backend

# or just the frontend
make frontend
```

The server will start on port 8081. You can then access the Web UI in your browser:

🔗 http://localhost:8081

## ⚙️ REST API Endpoints

The server exposes a set of endpoints for cluster management and raw data operations. All operations are performed against the currently active cluster.

### Health Check

| Method | Endpoint | Description                     |
| ------ | -------- | ------------------------------- |
| GET    | /health  | Check the server health status. |

### Cluster Management

| Method | Endpoint              | Description                                | Body Example                                        |
| ------ | --------------------- | ------------------------------------------ | --------------------------------------------------- |
| POST   | /api/clusters/connect | Connect to a new TiKV cluster.             | `{"pd_addrs": ["host:port"], "name": "production"}` |
| GET    | /api/clusters         | List all connected clusters.               | N/A                                                 |
| POST   | /api/clusters/switch  | Set an existing cluster as the active one. | `{"name": "production"}`                            |

### Raw KV Operations (Active Cluster)

| Method | Endpoint        | Description                            | Body Example                                       |
| ------ | --------------- | -------------------------------------- | -------------------------------------------------- |
| POST   | /api/raw/get    | Retrieve the value for a specific key. | `{"key": "mykey"}`                                 |
| POST   | /api/raw/put    | Insert or update a key-value pair.     | `{"key": "mykey", "value": "myvalue"}`             |
| POST   | /api/raw/delete | Delete a key-value pair.               | `{"key": "mykey"}`                                 |
| POST   | /api/raw/scan   | Scan a range of keys.                  | `{"start_key": "a", "end_key": "z", "limit": 100}` |

### Metrics

| Method | Endpoint | Description                            |
| ------ | -------- | -------------------------------------- |
| GET    | /metrics | PD and TiKV metrics from the instances |

## 🤝 Contributing

We welcome contributions from the community! If you're interested in making the TiKV Admin Web UI even better:

- Report bugs and suggest features.
- Submit pull requests.
- Join the discussion.
