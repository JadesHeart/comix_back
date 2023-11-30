async function postData(url, data) {
    try {
        const response = await fetch(url, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(data)
        });

        if (!response.ok) {
            throw new Error(`HTTP error! Status: ${response.status}`);
        }

        const reader = response.body.getReader();

        let receivedLength = 0;
        let chunks = [];

        while (true) {
            const { done, value } = await reader.read();

            if (done) {
                break;
            }

            chunks.push(value);
            receivedLength += value.length;
        }

        const blob = new Blob(chunks);

        const img = document.createElement('img');
        img.src = URL.createObjectURL(blob);

        document.body.appendChild(img);
    } catch (error) {
        console.error('Error:', error);
    }
}

const data = {

    "password": "LQz63rOj7F1LzN3uASRR",
    "tagName": "bigcock",
    "name": "Почтовая спешка"
};

postData('http://localhost:8082/files', data);
  