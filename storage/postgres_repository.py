"""PostgreSQL index for events and graph metadata."""

from __future__ import annotations


class PostgresRepository:
    def __init__(self, database_url: str) -> None:
        self.database_url = database_url
