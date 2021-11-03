import functools
import os

import click
import psycopg2

def connection(func):
    @functools.wraps(func)
    def connection_wrapper(*args, **kwargs):
        domain = os.environ["STORAGE_SERVICE_DATABASE_DOMAIN"]
        user = os.environ["STORAGE_SERVICE_DATABASE_USER"]
        password = os.environ["STORAGE_SERVICE_DATABASE_PASSWORD"]
        database_name = kwargs.pop(
            'database_name',
            os.environ["STORAGE_SERVICE_ADMIN_DATABASE"],
        )
        conn = None
        try:
            conn = psycopg2.connect(
                host=domain,
                dbname=database_name,
                user=user,
                password=password,
            )
            func(conn, *args, **kwargs)
        except Exception as err:
            click.echo(f"Database operation failed.")
            click.echo(f"exception: \'{str(err).strip()}\'")
        finally:
            if conn is not None:
                conn.close()
    return connection_wrapper
