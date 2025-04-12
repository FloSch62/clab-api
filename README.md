# Containerlab API Server

This project provides a standalone RESTful API server written in Go to interact with the [Containerlab](https://containerlab.dev/) command-line tool.

**WARNING:** This API allows executing commands (`clab`) on the host system, potentially with elevated privileges (`sudo`). It includes **highly insecure placeholder authentication**. Use with extreme caution, understand the security implications, and **do not expose this API publicly without significant hardening**, proper authentication (PAM, LDAP, etc.), authorization controls, and HTTPS.

## Features

*   **Deploy Labs:** Start labs from `.clab.yml` files located in user home directories.
*   **Destroy Labs:** Stop and clean up running labs.
*   **Inspect Labs:** Get details about running labs and nodes.
*   **List Labs:** View all running labs associated with the user.
*   **List Topologies:** List available `.clab.yml` files in the user's home directory.
*   **Authentication:** JWT-based (login required).
*   **Authorization:** API endpoints require a valid JWT.
*   **User Context:** Executes `clab` commands as the authenticated Linux user via `sudo`.
*   **Swagger/OpenAPI:** Auto-generated API documentation.

## Prerequisites

1.  **Go:** Version 1.20 or higher.
2.  **Containerlab:** The `clab` executable must be installed and available in the system's `PATH`.


## Setup

1.  **Clone the repository:**
    ```bash
    git clone <your-repo-url> clab-api
    cd clab-api
    ```
2.  **Install Go dependencies:**
    ```bash
    go mod tidy
    ```
3.  **Configure Environment:**
    Create a `.env` file in the project root (copy from `.env.example`) and **change `JWT_SECRET`** to a strong, random value.
    ```bash
    cp .env.example .env
    # Edit .env and set a strong JWT_SECRET
    ```
4.  **Generate Swagger Docs:**
    ```bash
    # Install swag if you haven't already
    # go install github.com/swaggo/swag/cmd/swag@latest
    swag init -g cmd/server/main.go
    ```
5.  **Build the server:**
    ```bash
    go build -o clab-api-server ./cmd/server/
    ```

## Running the Server

```bash
# Make sure you are running as the user configured in sudoers (e.g., 'apiuser' or 'root')
./clab-api-server
```
# Server Information

The server will start, typically on port **8080** (or as configured in `.env`).

---

## API Usage

### Authentication

**POST** `/login`  
JSON Body:
```json
{
  "username": "your_linux_user",
  "password": "your_linux_password"
}
```
- Returns a **JWT token**.  


### Authenticated Requests

Include the token in the `Authorization` header for all `/api/v1/*` requests:
Authorization: Bearer <your_jwt_token>

---

## Endpoints

- **Swagger UI:** http://localhost:8080/swagger/index.html

### Labs
- **POST** `/api/v1/labs`  
  Deploy a lab (requires `topologyFile` relative to user home).

- **GET** `/api/v1/labs`  
  List running labs for the user.

- **GET** `/api/v1/labs/{labName}`  
  Inspect a specific lab.

- **DELETE** `/api/v1/labs/{labName}`  
  Destroy a specific lab.

### Topologies
- **GET** `/api/v1/topologies`  
  List `.clab.yml` files in the user's home directory.

---

## Security Considerations (Reminder)

- **Sudo:** Requires careful, restrictive configuration. Running the API as root is common but increases risk.
- **Authentication:** The default credential check is insecure. Implement PAM or other robust methods.
- **Input Validation:** Sanitize all user inputs (filenames, lab names) to prevent injection attacks.
- **HTTPS:** Absolutely essential for any non-local deployment. Use a reverse proxy (Nginx, Caddy) to handle TLS termination.
- **Rate Limiting/Firewalling:** Protect the API from abuse.
- **Error Handling:** Avoid leaking sensitive information in error messages.

---

## Development

- Use `swag init -g cmd/server/main.go` to update Swagger docs after changing comments or models.
- Consider adding more robust error handling and logging.


---

**5. Build and Run**

1.  **Generate Swagger Docs:** `swag init -g cmd/server/main.go`
2.  **Build:** `go build -o clab-api-server ./cmd/server/`
3.  **Run:** `sudo ./clab-api-server` (Run as the user configured in `sudoers`, often root for simplicity in managing `sudo -u`).
4.  Access `http://localhost:8080/swagger/index.html` in your browser.
5.  Use `curl` or a tool like Postman/Insomnia to interact with the API, starting with `/login`. Remember to include the `Authorization: Bearer <token>` header for protected endpoints. Place your `.clab.yml` files in the home directory of the user you log in as (e.g., `/home/myuser/test.clab.yml`).