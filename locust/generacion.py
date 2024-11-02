import random
import os
import json

nombres = ["Monica", "Oliver", "Marco", "Mario", "Flor", "Josue", "Solis", "Luigi", "Luis", "Alejandra", "Abigail", "Pamela", "Maria", "Daniela", "Antonio"]
facultad = ["Ingenieria", "Agronomia"]

def generar_json():

    
    json_students = []
    for i in range(0,500):
        
        nombre = random.choice(nombres) 
        json_student = {
            "student": f"{nombre}_{i}",
            "age": random.randint(15,45),
            "faculty": random.choice(facultad),
            "discipline": random.randint(1,3)
        }

        json_students.append(json_student)
        
    with open("json_estudiantes.json", "+w") as file:
        json.dump(json_students, file, indent=4)


generar_json()

