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


@cli.command("create-database-user")
@click.argument("database_name")
@click.argument("username")
@click.password_option()
def create_database_user(database_name, username, password):
    database.create_user(username, password, database_name=database_name)


@cli.command("drop-database-user")
@click.argument("database_name")
@click.argument("username")
def drop_database_user(database_name, username):
    database.drop_user(username, database_name=database_name)


@cli.command("list-databases")
def list_databases():
    database.list_databases()


if __name__ == "__main__":
    cli()
