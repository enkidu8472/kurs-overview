
let gApiResponseAssetData = null;

async function fetchAssetDataAndRefreshChart() {

    const inputIsins = getElementsFromSelectList();
    const fromDate = document.getElementById("from-date").value;
    const toDate = document.getElementById("to-date").value;

    // stop function when no ISIN is selected
    if (inputIsins.length == 0) {
        return;
    }

    try {

        const response = await fetch('http://localhost:8080/api/getAssetData', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({
                isins: inputIsins,
                fromDate: fromDate,
                toDate: toDate
            })
        });

        if (!response.ok) {
            console.log(response)
            throw new Error('Failed to fetch fetch-AssetData:' + inputIsins);
        }

        gApiResponseAssetData = await response.json();
        console.log("i have the response !!!! 2")
        console.log(gApiResponseAssetData)
        gApiResponseAssetData = fixTimeFormatFromResponse(gApiResponseAssetData);
        updateChart1(inputIsins);
        updateChart2(inputIsins);
        createSummaryTable();

    } catch (error) {
        console.log(error)
        window.alert("In fetch-AssetData :: isin:" + inputIsins + "\n err37: " + error.message);
    }
}

let gApiResponsePossibleAssets = null;

async function fetchPossibleAssets() {

    try {

        const response = await fetch('http://localhost:8080/api/getPossibleAssets', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            }
        });

        if (!response.ok) {
            throw new Error('Failed to fetch getPossibleAssets');
        }

        gApiResponsePossibleAssets = await response.json();
        initSelectableAssets();
    } catch (error) {
        window.alert("In fetch-PossibleAssets err51: " + error.message);
    }
}


fetchPossibleAssets();

function fixTimeFormatFromResponse(responseData) {
    for (let i = 0; i < responseData.curves.length; i++) {
        let asset = responseData.curves[i];
        for (let j = 0; j < asset.values.length; j++) {
            let entry = asset.values[j];
            entry.date = new Date(entry.date);
        }
    }
    return responseData;
}


async function fetchAssetData4SelectWindow() {

    const inputIsins = gApiResponsePossibleAssets;
    const fromDate = "";
    const toDate = "";

    // stop function when no ISIN is selected
    if (inputIsins.length == 0) {
        return;
    }

    try {

        const response = await fetch('http://localhost:8080/api/getAssetData', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({
                isins: inputIsins,
                fromDate: fromDate,
                toDate: toDate
            })
        });

        if (!response.ok) {
            console.log(response)
            throw new Error('Failed to fetch fetch-AssetData:' + inputIsins);
        }

        selectWindowData = await response.json();
        console.log("i have the response 3");
        console.log(selectWindowData);
        selectWindowData = fixTimeFormatFromResponse(selectWindowData);
        createSelectionWindowTable(selectWindowData.curves);

    } catch (error) {
        console.log(error)
        window.alert("In fetch-AssetData :: isin:" + inputIsins + "\n err127: " + error.message);
    }
}