import click
import os

import psycopg2
from psycopg2.extensions import ISOLATION_LEVEL_AUTOCOMMIT

def create_database(database_name):
    domain = os.environ["STORAGE_SERVICE_DATABASE_DOMAIN"]
    user = os.environ["STORAGE_SERVICE_DATABASE_USER"]
    password = os.environ["STORAGE_SERVICE_DATABASE_PASSWORD"]
    admin_database = os.environ["STORAGE_SERVICE_ADMIN_DATABASE"]
    try:
        conn = psycopg2.connect(
            host=domain,
            dbname=admin_database,
            user=user,
            password=password,
        )
        conn.set_isolation_level(ISOLATION_LEVEL_AUTOCOMMIT) 
        cur = conn.cursor()
        cur.execute(f"CREATE DATABASE {database_name};")
        click.echo(f"database \"{database_name}\" created")
    except Exception as err:
        click.echo(f"database \"{database_name}\" could not be created")
        click.echo(f"exception: \'{str(err).strip()}\'")


def drop_database(database_name):
    domain = os.environ["STORAGE_SERVICE_DATABASE_DOMAIN"]
    user = os.environ["STORAGE_SERVICE_DATABASE_USER"]
    password = os.environ["STORAGE_SERVICE_DATABASE_PASSWORD"]
    admin_database = os.environ["STORAGE_SERVICE_ADMIN_DATABASE"]
    try:
        conn = psycopg2.connect(
            host=domain,
            dbname=admin_database,
            user=user,
            password=password,
        )
        conn.set_isolation_level(ISOLATION_LEVEL_AUTOCOMMIT) 
        cur = conn.cursor()
        cur.execute(f"DROP DATABASE {database_name};")
        click.echo(f"database \"{database_name}\" dropped")
    except Exception as err:
        click.echo(f"database \"{database_name}\" could not be dropped")
        click.echo(f"exception: \'{str(err).strip()}\'")
