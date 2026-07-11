#!/usr/bin/env node

// This handles the plugin event. In this example, we are handling the "beforeTransform" event. 
// Available events are: "beforeTransform" and "afterTransform". You can add your own logic here to modify the data based on the event.
function handle(version, event, data) {
    if(event === "beforeTransform") {
        if(!data.input) data.input = {};
        data.input.value = "This field was added by the beforeTransform plugin.";
    }
    return data;
}

function main() {
    
    process.stdin.on('data', function (data) {
        const line = data.trim();
        let response = {}
        
        if(!line) {
            return;
        }

        try {
            // The plugin receives a JSON payload from watt. 
            // The payload contains the version, event and data.
            // Data contains the environment variables ( env ), the input data ( input ) and the result of the transformation ( result ).
            const request = JSON.parse(line);
            const version = request.version || "";
            const event = request.event || "";
            const data = request.data || {};
            const updatedData = handle(version, event, data);

            response = {
                status: "success",
                data: updatedData
            };

        } catch (error) {
            response = {
                status: "error",
                error: error.message
            };
        }

        process.stdout.write(JSON.stringify(response) + "\n");
        process.stdout.flush();
        process.exit(0);
    });
}

main();