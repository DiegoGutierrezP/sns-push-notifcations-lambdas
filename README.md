# Digital push notifications lambdas

This project contains a set of AWS Lambda functions designed to manage push notification workflows, including device registration, topic subscriptions, message publishing, and delivery status logging.

Lambdas Included:

- Register Device – Stores device metadata and creates SNS endpoints.
- Update Device – Updates device information such as tokens or OS data.
- Device Subscribe – Subscribes a device to a specific topic.
- Device Unsubscribe – Removes a device from a topic.
- Publish Notification – Sends push notifications to devices or topics.
- Delivery Status Log Processor – Processes SNS delivery status logs and stores them for auditing.

## Local Development & Testing

1. Build the docker image

```
docker build -t lmbd-digital-push-notifications .
```

2. Start the local environment

```
docker-compose up -d
```

3. Rebuild a specific Lambda when making changes

Replace <folder> with the folder name inside /cmd.

```
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-s -w" -o ./bin/<folder> ./cmd/<folder>
```
