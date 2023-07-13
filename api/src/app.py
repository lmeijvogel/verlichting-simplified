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

ALLOWED_SCENES = ["scene.uit", "scene.ochtend", "scene.middag", "scene.avond", "scene.nacht"];
ALLOWED_SWITCHES = ["switch.elektrische_deken", "switch.mechanische_ventilatie", "switch.tv_meubel"];

@app.route("/api/scenes", methods=['GET'])
def scenes():
    url = "http://%s/api/states" % API_HOST

    headers = {"Authorization": "Bearer %s" % API_TOKEN}

    response = requests.get(url, headers=headers)
    content = response.content

    entities = json.loads(content)

    scenes = [scene_to_json(entity) for entity in entities if entity["entity_id"] in ALLOWED_SCENES]
    return {
            "scenes": scenes,
            "fullResponse": entities
    }

@app.route("/api/start_scene/<scene_name>", methods=['POST'])
def start_scene(scene_name):
    if (scene_name in ALLOWED_SCENES):
        url = "http://%s/api/services/scene/turn_on" % API_HOST

        headers = {"Authorization": "Bearer %s" % API_TOKEN}

        data = {"entity_id": scene_name }

        response = requests.post(url, headers=headers, data=json.dumps(data))

        updated_scene = json.loads(response.content)[0]

        return scene_to_json(updated_scene)
    else:
        return "Scene not found", 404

@app.route("/api/switches", methods=['GET'])
def switches():
    url = "http://%s/api/states" % API_HOST

    headers = {"Authorization": "Bearer %s" % API_TOKEN}

    response = requests.get(url, headers=headers)
    content = response.content

    entities = json.loads(content)

    return {
            "switches": [switch_to_json(entity) for entity in entities if entity["entity_id"] in ALLOWED_SWITCHES],
            "fullResponse": entities
    }

@app.route("/api/set_switch/<switch_id>/<switch_state>", methods=['POST'])
def set_switch(switch_id, switch_state):
    if (switch_id in ALLOWED_SWITCHES):

        command = "turn_on" if switch_state == "on" else "turn_off"

        url = "http://%s/api/services/switch/%s" % [API_HOST, command]

        headers = {"Authorization": "Bearer %s" % API_TOKEN}

        data = {"entity_id": switch_id }

        try: 
            response = requests.post(url, headers=headers, data=json.dumps(data))

            updated_switch = json.loads(response.content)[0]

            return switch_to_json(updated_switch)
        except json.decoder.JSONDecodeError as exc:
            print(response)

def scene_to_json(scene):
    return {
            "id": scene["entity_id"],
            "friendlyName": scene["attributes"]["friendly_name"],
            "lastActivated": scene["state"]
            }

def switch_to_json(scene):
    return {
            "id": scene["entity_id"],
            "friendlyName": scene["attributes"]["friendly_name"],
            "state": scene["state"]
            }

if __name__ == "__main__":
    app.run(host='0.0.0.0', port=3000)
