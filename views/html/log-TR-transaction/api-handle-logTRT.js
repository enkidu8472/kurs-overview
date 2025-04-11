

let gIsinDownloadResponse = null;

async function pushInfoToBackend(inputString) {
  
    try {

        const response = await fetch('http://localhost:8080/api/putTRtransaction', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({
                clipboard: inputString
            })
        });

        if (!response.ok) {
            throw new Error('Failed to fetch fetch-AssetData:' + inputString);
        }

        gIsinDownloadResponse = await response.json();

    } catch (error) {
        window.alert("In fetch-AssetData :: isin:" + inputString + "\n err27: " + error.message);
    }
}

