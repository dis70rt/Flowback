# Flowback Backend Documentation

Welcome to the Flowback Backend. This guide will show you how to start the local development environment, run the database migrations, and expose your local server to the internet so Razorpay can send you webhooks.

## 1. Start the Infrastructure (Database & Redis)
We use Docker to run our dependencies. Make sure Docker is running on your machine, then execute:

```bash
# From the root of the project
docker compose up -d
```
This will start:
- **Postgres (v16)**: For our application data and immutable audit logs.
- **Redis (v7)**: For our Asynq task queues and Pub/Sub event routing.

## 2. Start the API Server
Our Go API handles incoming webhooks and serves the frontend dashboard. 

*(Note: The API automatically runs `goose` database migrations via `go:embed` the moment it starts up. You do not need to run migrations manually!)*

```bash
# From the apps/backend directory
go run cmd/api/main.go
```
You should see: `STARTING: Flowback Backend listening on port 8080...`

## 3. Expose Localhost to the Internet (Cloudflare Tunnel)
Razorpay needs a public HTTPS URL to send webhooks to. Since your backend is running locally on port 8080, we use `cloudflared` to securely tunnel it to the web.

Open a **new terminal tab** and run:

```bash
cloudflared tunnel --url http://localhost:8080
```

*(Note: If you are using the custom domain from our previous setup, run your persistent tunnel mapping to `flowback-webhook.saikat.in`)*

Cloudflare will give you a public URL (e.g., `https://random-words.trycloudflare.com`). **Copy this URL.**

## 4. Configure Razorpay Webhooks
1. Log in to the [Razorpay Dashboard](https://dashboard.razorpay.com/) (Ensure you are in **Test Mode**).
2. Go to **Account & Settings -> Webhooks -> Add New Webhook**.
3. **Webhook URL:** Paste your Cloudflare URL and append the route: 
   `https://<your-cloudflare-url>/webhooks/razorpay`
4. **Secret:** Enter the secret defined in your `.env` file (`RAZORPAY_WEBHOOK_SECRET`).
5. **Events to Subscribe to:**
   - `subscription.pending`
   - `subscription.halted`
   - `payment.failed`
   - `payment_link.paid`

## 5. Next Steps: Background Workers
Once the API is running and receiving webhooks, you will need to start the AI Background Workers in separate terminal tabs to process the queues (Soft Declines, Hard Declines, etc.).

```bash
# Example: Start the Recovery Worker (Coming Soon)
go run cmd/worker/main.go
```
