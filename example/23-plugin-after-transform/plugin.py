#!/usr/bin/env python3
import json
import os
import sys

def handle(version, event, data):
    if event == "afterTransform":
        data["result"]["resource"]["additional"] = {"value": "set by plugin"}
    return data

def main():
    line = sys.stdin.readline()

    try:
        req = json.loads(line)
        version = req.get("version", "")
        data = req.get("data", {})
        event = req.get("event", "")

        if event == "afterTransform":
            updatedData = handle(version, event, data)

        response = {"status": "success", "data": updatedData}
    except Exception as e:
        response = {"status": "error", "error": str(e)}

    sys.stdout.write(json.dumps(response) + "\n")
    sys.stdout.flush()
    os._exit(0)


if __name__ == "__main__":
    main()