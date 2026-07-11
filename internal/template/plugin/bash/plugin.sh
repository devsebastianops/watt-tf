#!/bin/bash

# This handles the plugin event. In this example, we are handling the "beforeTransform" event. 
# Available events are: "beforeTransform" and "afterTransform". You can add your own logic here to modify the data based on the event.
function handle(version, event, data) {
    if [[ "$event" == "beforeTransform" ]]; then    
        data=$(echo "$data" | jq '.input.value = "This value was modified by the plugin before the transformation."')
    fi

    # Return the modified data
    echo "$data"
}

function main() {
    # The plugin receives a JSON payload from watt. 
    # The payload contains the version, event and data.
    # Data contains the environment variables ( env ), the input data ( input ) and the result of the transformation ( result ).
    stdin=$(cat -)
    
    req=$(echo "$stdin" | jq -r '.')
    version=$(echo "$req" | jq -r '.version')
    event=$(echo "$req" | jq -r '.event')
    data=$(echo "$req" | jq -r '.data')

    response=$(handle "$version" "$event" "$data")
    
    echo "$response"
}

main