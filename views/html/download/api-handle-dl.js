

let gIsinDownloadResponse = null;

async function downloadIsinForBackend(inputIsin) {
  
    try {

        const response = await fetch('http://localhost:8080/api/downloadFullIsinData', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({
                isin: inputIsin
            })
        });

        if (!response.ok) {
            throw new Error('Failed to fetch fetch-AssetData:' + inputIsin);
        }

        gIsinDownloadResponse = await response.json();

    } catch (error) {
        window.alert("In fetch-AssetData :: isin:" + inputIsin + "\n err27: " + error.message);
    }
}

