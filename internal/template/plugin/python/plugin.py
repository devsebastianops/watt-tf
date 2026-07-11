#!/usr/bin/env python3
import json
import os
import sys

# This handles the plugin event. In this example, we are handling the "beforeTransform" event. 
# Available events are: "beforeTransform" and "afterTransform". You can add your own logic here to modify the data based on the event.
def handle(version, event, data):
    if event == "beforeTransform":
        data["input"]["value"] = "This field was added by the beforeTransform plugin."
    return data

def main():
    line = sys.stdin.readline()
        
    try:
        # The plugin receives a JSON payload from watt. 
        # The payload contains the version, event and data.
        # Data contains the environment variables ( env ), the input data ( input ) and the result of the transformation ( result ).
        req = json.loads(line)
        version = req.get("version", "")
        data = req.get("data", {})
        event = req.get("event", "")

        updatedData = handle(version, event, data)

        # watt expects a JSON response with a status ( "success" or "error" ) and the complete data object with its updates made by the plugin.
        response = {"status": "success", "data": updatedData}
    except Exception as e:
        response = {"status": "error", "error": str(e)}

    sys.stdout.write(json.dumps(response) + "\n")
    sys.stdout.flush()
    os._exit(0)


if __name__ == "__main__":
    main()