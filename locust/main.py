from locust import HttpUser, TaskSet, task, between
import json
import random

with open("json_estudiantes.json", "r") as file:
    datos = json.load(file)

endpoints = ["Agronomia", "Ingenieria"]

class StudentBehavior(TaskSet):
    @task
    def post_request(self):

        estudiante = random.choice(datos)
        endpoint = random.choice(endpoints)

        url = f"http://35.233.190.115.nip.io/{endpoint}"
        headers = {'Content-Type': 'application/json'}

        self.client.post(url, headers=headers, data=json.dumps(estudiante))
        

class WebsiteUser(HttpUser):
    tasks = [StudentBehavior]
    wait_time = between(1,3)