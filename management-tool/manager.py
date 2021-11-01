import click
import os

from commands import database

@click.group()
def cli():
    pass


@cli.command("create-database")
@click.argument("database_name")
def create_database_handler(database_name):
    database.create_database(database_name)


@cli.command("drop-database")
@click.argument("database_name")
def drop_database_handler(database_name):
    database.drop_database(database_name)


if __name__ == "__main__":
    cli()
