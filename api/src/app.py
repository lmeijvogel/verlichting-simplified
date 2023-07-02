from dotenv import load_dotenv
from flask import Flask
import os
import sys

import requests

import json

load_dotenv()

app = Flask(__name__)

API_HOST = os.getenv("API_HOST")
API_TOKEN = os.getenv("API_TOKEN")

@app.route("/api/scenes", methods=['GET', 'POST'])
def scenes():
    url = "http://%s/api/states" % API_HOST

    headers = {"Authorization": "Bearer %s" % API_TOKEN}

    response = requests.get(url, headers=headers)
    content = response.content

    entities = json.loads(content)

    return {
            "scenes": [scene_to_json(entity) for entity in entities if (entity["entity_id"].startswith('scene.'))],
            "fullResponse": entities
    }

def scene_to_json(scene):
    return {
            "id": scene["entity_id"],
            "friendlyName": scene["attributes"]["friendly_name"],
            "lastActivated": scene["state"]
            }

@app.route("/api/start_scene/<scene_name>", methods=['POST'])
def start_scene(scene_name):
    url = "http://%s/api/services/scene/turn_on" % API_HOST

    headers = {"Authorization": "Bearer %s" % API_TOKEN}

    data = {"entity_id": scene_name }

    response = requests.post(url, headers=headers, data=json.dumps(data))

    updated_scene = json.loads(response.content)[0]

    return scene_to_json(updated_scene)

if __name__ == "__main__":
    app.run(host='0.0.0.0', port=3000)
