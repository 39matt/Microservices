FROM python:3.13-slim as BUILDER

WORKDIR /app

COPY requirements.txt .

RUN --mount=type=cache,target=/root/.cache/pip \
    pip3 install -r requirements.txt

COPY . .

EXPOSE 8100

CMD ["python3", "main.py"]