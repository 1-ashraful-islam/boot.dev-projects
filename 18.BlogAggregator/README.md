# Blog Aggregator in Go

This project is a blog aggregator service in Go. It utilizes RESTful API that fetches data from remote locations and stores them in a production-ready database tools like PostgreSQL, SQLc, Goose, and pgAdmin.

It also utilizes a long-running service worker that reaches out over the internet to fetch data from remote locations.

## Usage

1. Rename the `.env.example` file to `.env` and fill in the required environment variables.

```bash
cp .env.example .env
```

2. Run the following command to start postgresql server using docker-compose

  ```bash
  docker-compose up -d
  ```

3. Run the following command to stop postgresql server using docker-compose

  ```bash
  docker-compose down
  ```
