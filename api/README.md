# OSS Health API

A FastAPI-based backend service for the OSS Health project, providing GraphQL API endpoints and GitHub OAuth authentication.

## Features

- FastAPI-based REST and GraphQL API
- GitHub OAuth authentication
- SQLAlchemy with async PostgreSQL support
- Alembic database migrations
- Comprehensive test suite with pytest
- Docker support for containerization

## Prerequisites

- Python 3.12 or higher
- Poetry for dependency management
- PostgreSQL database
- Docker (optional)

## Installation

1. Clone the repository:
```bash
git clone <repository-url>
cd api
```

2. Install dependencies using Poetry:
```bash
poetry install
```

3. Set up your environment variables (create a `.env` file):
```env
DATABASE_URL=postgresql+asyncpg://user:password@localhost:5432/dbname
GITHUB_CLIENT_ID=your_github_client_id
GITHUB_CLIENT_SECRET=your_github_client_secret
```

## Development

1. Activate the virtual environment:
```bash
poetry shell
```

2. Run database migrations:
```bash
alembic upgrade head
```

3. Start the development server:
```bash
uvicorn main:app --reload
```

The API will be available at `http://localhost:8000`

## API Documentation

- REST API documentation: `http://localhost:8000/docs`
- GraphQL playground: `http://localhost:8000/graphql`

## Testing

Run the test suite:
```bash
pytest
```

## Docker

Build and run the Docker container:
```bash
docker build -t oss-health-api .
docker run -p 8000:8000 oss-health-api
```

## Project Structure
