from __future__ import annotations

from datetime import datetime
from airflow.decorators import dag, task

@dag(
    dag_id="hello_dag",
    start_date=datetime(2025, 1, 1),
    schedule="@daily",
    catchup=False,
    tags=["demo"],
)
def hello_dag():
    @task
    def say_hello():
        print("Hello from Airflow!")

    @task
    def say_bye():
        print("Bye!")

    say_hello() >> say_bye()

dag = hello_dag()