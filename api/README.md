# OSS Health API

A FastAPI-based backend service for the OSS Health project, providing GraphQL API endpoints and GitHub OAuth authentication.

## Features

- FastAPI-based REST and GraphQL API
- GitHub OAuth authentication
- SQLAlchemy with async PostgreSQL support
- Alembic database migrations
- Comprehensive test suite with pytest
- Docker support for containerization
- Modern Python tooling with Poetry and Ruff

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
poetry run uvicorn main:app --reload
```

The API will be available at `http://localhost:8000`

## Code Quality

The project uses Ruff for linting and code formatting. Ruff is configured to:
- Check for errors (E)
- Check for code style (F)
- Sort imports (I)

To run the linter:
```bash
poetry run ruff check .
```

To automatically fix issues:
```bash
poetry run ruff check --fix .
```

## API Documentation

- REST API documentation: `http://localhost:8000/docs`
- GraphQL playground: `http://localhost:8000/graphql`

## Testing

Run the test suite:
```bash
poetry run pytest
```

## Docker

Build and run the Docker container:
```bash
docker build -t oss-health-api .
docker run -p 8000:8000 oss-health-api
```

## Dependencies

Main dependencies:
- FastAPI
- Strawberry GraphQL
- SQLAlchemy
- Alembic
- asyncpg
- Python-Jose (for JWT)
- PyYAML

Development dependencies:
- Ruff (linting and formatting)
- Pytest (testing)
- MyPy (type checking)

## Contributing

1. Fork the repository
2. Create a feature branch
3. Commit your changes
4. Push to the branch
5. Create a Pull Request

Before submitting a PR, ensure your code passes all checks:
```bash
poetry run ruff check .
poetry run pytest
```

## License

This project is licensed under the GNU Affero General Public License v3.0 - see the [LICENSE](../LICENSE) file for details.

The GNU Affero General Public License is a free, copyleft license for software and other kinds of works, specifically designed to ensure cooperation with the community in the case of network server software.

## Author

JP Wienekus (jpwienekus@gmail.com)
