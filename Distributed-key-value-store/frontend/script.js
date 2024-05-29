document.addEventListener("DOMContentLoaded", function () {
    // API base URL
    const baseUrl = "http://localhost:9090";

    // Get form elements
    const setForm = document.getElementById("set-form");
    const getKeyForm = document.getElementById("get-form");
    const nodeList = document.getElementById("node-list");

    // Function to fetch node status
    function fetchNodeStatus() {
        fetch(baseUrl + "/status")
            .then(response => response.json())
            .then(data => {
                // Clear existing node list
                nodeList.innerHTML = "";

                // Display node status
                data.nodes.forEach(node => {
                    const li = document.createElement("li");
                    li.textContent = `${node.id}: ${node.status}`;
                    nodeList.appendChild(li);
                });
            })
            .catch(error => console.error("Error fetching node status:", error));
    }

    // Function to handle setting key-value pairs
    setForm.addEventListener("submit", function (event) {
        event.preventDefault();
        const formData = new FormData(setForm);
        const key = formData.get("key");
        const value = formData.get("value");

        // Make API request to set key-value pair
        fetch(baseUrl + "/set", {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify({ key: key, value: value })
        })
            .then(response => {
                if (response.ok) {
                    console.log("Key-value pair set successfully.");
                    fetchNodeStatus(); // Refresh node status
                } else {
                    console.error("Failed to set key-value pair.");
                }
            })
            .catch(error => console.error("Error setting key-value pair:", error));
    });

    // Function to handle getting value by key
    getKeyForm.addEventListener("submit", function (event) {
        event.preventDefault();
        const formData = new FormData(getKeyForm);
        const key = formData.get("get-key");

        // Make API request to get value by key
        fetch(baseUrl + `/get?key=${key}`)
            .then(response => response.json())
            .then(data => {
                if (data.value) {
                    console.log(`Value for key "${key}": ${data.value}`);
                } else {
                    console.log(`Key "${key}" not found.`);
                }
            })
            .catch(error => console.error(`Error getting value for key "${key}":`, error));
    });

    // Initial fetch for node status
    fetchNodeStatus();
});
