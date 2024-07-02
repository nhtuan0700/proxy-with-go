## What is the repository used for?
It is a basic setup of Proxy, Reverse Proxy, Load Balancing

# Explain:
## Proxy:
- Proxy server acts as an intermediary between clients (such as web browsers) and other servers.
  - key functions: 
    - Anonymity and Security: Proxy servers can hide a client's IP address, providing anonymity, and can also filter web traffic to block malicious content.
    - Content Filtering: Organizations often use proxies to control access to certain websites or content, enforcing policies on web usage.
    - Performance: Proxies can cache frequently requested resources, improving speed and reducing bandwidth usage.
#### How to run:
1. Set `.env` with:
    ```env
    HTTP_PROXY="server.local:18080"
    ```
    > NOTE: Don't use localhost, otherwise, proxy is not working

2. Run proxy server
    ```sh
    cd proxy && go run .
    ```

3. Run client to test it
    ```sh
    cd client && go run .
    ```
    - Make sure that log is not show "Can not find proxy url"

## Reverse Proxy + Load Balancing
- Reverse proxy: on the other hand, sits between clients and servers, but it primarily serves to manage traffic on behalf of servers.
    - key functions: 
      - Load Balancing: Distributing client requests across multiple servers to optimize resource utilization and improve reliability.
      - SSL Termination: Handling SSL/TLS encryption and decryption, offloading this task from backend servers to improve performance.
      - Content Caching: Storing copies of frequently accessed resources closer to clients, reducing latency and server load.
    > In this section, I implement only load balancing

- Load balancing: Distributing client requests across multiple servers to optimize resource utilization and improve reliability.

#### What have i done?
- The reverse proxy implements load balancer which is distributing 3 servers
- Each server will be rotated to handle request
- After 10s, the load balancer will ping to each server to verify if the host is available
- If within that 10s, there is a server that is not available, the load balancer will transfer to another one. If there is no available server, it will return a response "Service unavailable"

#### how to run: 
1. Set url for reverse proxy server, set in `.env`
    ```env
    REVERSE_PROXY_URL=http://localhost:28080
    ```
2. Assume there are 3 hosts available, please set `.env` for them:
    ```env
    HOST_1=http://localhost:38080
    HOST_2=http://localhost:38081
    HOST_3=http://localhost:38082
    ```
3. Run reverse proxy server
    ```sh
    cmd reverse_proxy && go run .
    ```
4. Run client to test it
    ```sh
    cd client && go run .
    ```
