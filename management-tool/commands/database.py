import click
import os

import psycopg2
from psycopg2.extensions import ISOLATION_LEVEL_AUTOCOMMIT

from commands.decorators import connection

@connection
def create_database(conn, database_name):
    conn.set_isolation_level(ISOLATION_LEVEL_AUTOCOMMIT) 
    cur = conn.cursor()
    cur.execute(f"CREATE DATABASE {database_name};")
    click.echo(f"database \"{database_name}\" created")


@connection
def drop_database(conn, database_name):
    conn.set_isolation_level(ISOLATION_LEVEL_AUTOCOMMIT) 
    cur = conn.cursor()
    cur.execute(f"DROP DATABASE {database_name};")
    click.echo(f"database \"{database_name}\" dropped")
