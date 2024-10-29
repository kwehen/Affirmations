async function generateAffirmation() {
    const spinner = document.getElementById('spinner');
    const responseDiv = document.getElementById('response');

    // Show spinner
    spinner.style.display = 'block';
    responseDiv.innerHTML = ''; // Clear previous response

    try {
        const response = await fetch('/api/chat', {
            method: 'POST',
            headers: {
                'Content-Type': 'text/plain'
            },
            body: JSON.stringify({})
        });

        if (!response.ok) {
            throw new Error(response.statusText);
        }

        const data = await response.json();
        responseDiv.innerHTML = data.message || 'No message received';
    } catch (error) {
        responseDiv.innerHTML = 'Error: ' + error.message;
    } finally {
        // Hide spinner
        spinner.style.display = 'none';
    }
}
