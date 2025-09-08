# unnamed location rating app

unnamed location rating app is a social discovery app for exploring and ranking cities, neighborhoods, and areas through your own experiences. By writing about experiences in different places and highlighting hidden gems, each rating is ranked relative to past reviews to create a personalized list of top destinations. Share stories, drop recommendations, and see how your opinions compare with your friends'!

<img width="1728" height="971" alt="image" src="https://github.com/user-attachments/assets/d5993fa5-a268-4b6f-81b5-9bf795e4327d" />

## Development

File structure:
```
pin/
├── backend/           # Golang backend service
├── mobile/            # React Native frontend
├── .gitignore
└── README.md
```

### Run Locally

Backend:
```bash
cd backend
go run ./cmd/server
```

Frontend:
```bash
cd mobile
npm start
```

### Database (Postgres via Docker)

Prereqs: Docker Desktop

Start Postgres:
```bash
docker compose up -d db
```

Run migrations (uses `migrate/migrate` container):
```bash
cd backend
make migrate-up
```

Create a new migration:
```bash
cd backend
make migrate-create name=add_users_table
```

Environment variables (defaults shown):
```bash
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=pin
DB_HOST=localhost
DB_PORT=5432
DATABASE_URL=postgres://$DB_USER:$DB_PASSWORD@$DB_HOST:$DB_PORT/$DB_NAME?sslmode=disable
```
