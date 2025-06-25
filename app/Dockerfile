FROM python:3.13.5

RUN pip install poetry

WORKDIR /app

COPY pyproject.toml poetry.lock* ./

RUN poetry config virtualenvs.create false \
  && poetry install --no-root --only main

COPY . .

CMD ["poetry", "run", "uvicorn", "api.main:app", "--host", "0.0.0.0"]
