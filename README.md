# Summary
This project aims to provide a basic example of how to create an application that monitors certain elements using OpenTelemetry, storing those metrics in **SigNoz**, an open-source software for observability and performance monitoring.

To achieve this, three endpoints have been created, introducing random latencies through sleep and requests that fail 30% of the time to visualize the data in **SigNoz**.

# Prerequisites
## Certificate
Since we will be working with (SigNoz Cloud)[https://signoz.io/], we need to download the certificate corresponding to our region, in this case, Europe.
```
openssl s_client -showcerts -connect ingest.eu.signoz.cloud:443 </dev/null 2>/dev/null | openssl x509 -outform PEM > ca-cert.pem
``` 


To view different regions, you can visit the [SigNoz website](https://signoz.io/docs/ingestion/signoz-cloud/overview/).

## Ingestion Key
To ingest data into SigNoz, we need an **Ingestion Key** if using their Cloud service.
To obtain this, navigate to *Settings -> Ingestion Key*.
![Ingestion Key](./images/ingestion_key.png)

## Environment Variables
Before running the application, ensure you have set the following environment variables:
- `OTEL_EXPORTER_OTLP_ENDPOINT=ingest.eu.signoz.cloud`
- `SIGNOZ_ACCESS_TOKEN=<YOUR_INGESTION_KEY>`
- `CA_CERT_PATH=<PATH_TO_CERT>`

### Example
```
export OTEL_EXPORTER_OTLP_ENDPOINT=ingest.eu.signoz.cloud
export SIGNOZ_ACCESS_TOKEN=<YOUR_INGESTION_KEY>
export CA_CERT_PATH=ca_cert.pem
```

# Test
Now you can test the application using the following endpoints:
- `/process`
- `/add-to-cart`
- `/remove-from-cart`
```
> curl localhost:8080/process
Request processed successfully
> curl localhost:8080/process
Error processing the request

```