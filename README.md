# cogniboard
docs at: http://127.0.0.1:8888/v1/api/docs

## Development
1. Install devbox
2. Run `devbox shell`
3. `docker compose up -d`
4. `SERVICE_POSTGRES_DSN="postgres://cogniboard:password@localhost:5432/cogniboard?sslmode=disable" go run cmd/http/main.go`
   1. Running the command will automatically generate the `openapi3.yaml` file.
5. Go to web, run `bunx orval` to generate the react query hooks