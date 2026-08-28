# Flowback

Flowback is an audit microservice that ingests Razorpay webhooks, classifies payment decline events, and broadcasts decisions to downstream services.

## Architecture

1. **API**: Receives Razorpay webhooks, verifies HMAC signatures, and enqueues the raw payload to a Redis queue via Asynq.
2. **Worker**: Consumes tasks from Asynq, parses the JSON payload, classifies the decline type (Hard, Soft, Unknown), and broadcasts the decision.
3. **PubSub**: The worker publishes decisions to Redis PubSub channels (`decline:hard`, `decline:soft`, or `decline:unknown`).

## Prerequisites

* Go 1.22 or higher
* Docker and Docker Compose (for Redis and PostgreSQL)

## Setup

1. Copy `.env.example` to `.env` in the root directory and configure your variables.
2. Configure your Cloudflare Tunnel:
   * Go to the Cloudflare Zero Trust Dashboard -> Networks -> Tunnels.
   * Create a new tunnel (select Cloudflared) and copy the generated token.
   * Add this token to your `.env` file:
     ```env
     CLOUDFLARE_TUNNEL_TOKEN=your_token_here
     ```
   * In the tunnel configuration, click the **Published application routes** tab (do not use Hostname routes).
   * Click **Add a published application route**.
   * Set your desired Subdomain and select your Domain. Leave Path empty.
   * Under Service, set Type to `HTTP` and URL to exactly `host.docker.internal:8080`.
   * Save the route. *(Note: If you receive a "DNS record already exists" error, either choose a different subdomain or delete the conflicting record from your main Cloudflare DNS settings).*
3. Start the infrastructure:
   ```bash
   docker compose up -d
   ```

## Running the Services

The application is split into two separate processes. Run them in separate terminals.

### 1. Start the API Server
```bash
cd apps/backend
go run cmd/api/main.go
```
The API will listen on `http://localhost:8080`.

### 2. Start the Worker
```bash
cd apps/backend
go run cmd/worker/main.go
```
The worker will connect to Redis and begin polling for tasks.

## Testing End to End

You can trigger a real webhook flow by using the provided seeder. The seeder calls the Razorpay API to intentionally fail a transaction, which prompts Razorpay to send a `payment.failed` webhook to your configured webhook URL.

1. Ensure your `.env` file contains valid Razorpay API keys:
   ```env
   RAZORPAY_KEY_ID=your_key_id
   RAZORPAY_KEY_SECRET=your_key_secret
   ```
2. Ensure both the API and Worker are running.
3. Run the seeder tool:
   ```bash
   cd apps/backend
   go run cmd/seeder/main.go -task=card
   ```
4. Check the Worker terminal. Once Razorpay sends the webhook to your API, the worker will consume the task, classify the decline type (e.g., Soft decline for insufficient funds), and broadcast the decision to the PubSub channel.
