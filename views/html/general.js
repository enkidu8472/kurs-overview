
//import { downloadIsinForBackend } from './download/api-handle-dl.js';

// reset at startup
document.getElementById("from-date").value = "";
document.getElementById("to-date").value = "";



function initSelectableAssets() {

    // > load options to be able to select from 
    //const possibleAssets = ['Asset 1', 'Asset 2', 'Asset 3']

    // put by default the first element from the possible assets to the Selected Asseted list
    const select1Element = document.getElementById("selected-assets");
    const optionA = document.createElement("option");
    optionA.textContent = gApiResponsePossibleAssets[0];
    select1Element.appendChild(optionA);

    // all the other go to Available Assets list
    const select2Element = document.getElementById("available-assets");
    for (let i = 1; i < gApiResponsePossibleAssets.length; i++) {
        const optionB = document.createElement("option");
        optionB.textContent = gApiResponsePossibleAssets[i];
        select2Element.appendChild(optionB);
    }

    // inital drawing of the first graph
    fetchAssetDataAndRefreshChart();

}

function getElementsFromSelectList() {

    const selectElementOptions = document.getElementById("selected-assets").options;
    const optionsArray = [...selectElementOptions];      // Convert the options collection to an array using the spread operator    
    var valuesArray = optionsArray.map(function (option) {
        return option.value;
    });
    return valuesArray;
}

function createSummaryTable() {
    console.log("createSummaryTable 1.7")

    const tableContainer = document.getElementById("summaryTable");
    const summaryDuration = document.getElementById("summaryDuration");
    summaryDuration.innerHTML = gApiResponseAssetData.PeriodeDuration.Days + " (" +  gApiResponseAssetData.PeriodeDuration.Years + ")";
    tableContainer.innerHTML = "";

    // Isin
    let newRow = document.createElement("tr"); // Create a new row element
    for (let i = 0; i < gApiResponseAssetData.curves.length; i++) {

        const tableCell = document.createElement("td");
        const colorLabel = document.createElement("div");

        colorLabel.classList.add("lable-summary-table");
        colorLabel.style.backgroundColor = gApiResponseAssetData.curves[i].color;
        tableCell.appendChild(colorLabel)

        const textNode = document.createTextNode(gApiResponseAssetData.curves[i].name);
        tableCell.appendChild(textNode);
        tableCell.setAttribute("title", "ISIN");
        newRow.appendChild(tableCell);
        tableCell.onclick = function(){create_data_overview_right(gApiResponseAssetData.curves[i].name)};
    }
    tableContainer.appendChild(newRow);

    // percentagePerYear
    newRow = document.createElement("tr"); // Create a new row element
    for (let i = 0; i < gApiResponseAssetData.curves.length; i++) {
   
        let name = gApiResponseAssetData.curves[i].additionalInfo.nickname 
        const title = gApiResponseAssetData.curves[i].additionalInfo.title
        if (name === "" || name === undefined) {
            name = title          
            if (name.length > 16) {
                name = name.substring(0, 14) + "..."
            }
        }

        const tableCell = document.createElement("td");
        const textNode = document.createTextNode(name);
        tableCell.appendChild(textNode);
        
        tableCell.setAttribute("title", "Nickname or title of asset.\nTitle: " + title);
        newRow.appendChild(tableCell);
    }
    tableContainer.appendChild(newRow);

    // percentagePerYear
    newRow = document.createElement("tr"); // Create a new row element
    for (let i = 0; i < gApiResponseAssetData.curves.length; i++) {

        const tableCell = document.createElement("td");
        const textNode = document.createTextNode(gApiResponseAssetData.curves[i].percentagePerYear.toFixed(1));
        tableCell.appendChild(textNode);
        
        tableCell.setAttribute("title", "percentage change over current period by year, base is first value from periode.");
        newRow.appendChild(tableCell);
    }
    tableContainer.appendChild(newRow);

    // totalPercentage
    newRow = document.createElement("tr"); // Create a new row element
    for (let i = 0; i < gApiResponseAssetData.curves.length; i++) {

        const tableCell = document.createElement("td");
        const textNode = document.createTextNode(gApiResponseAssetData.curves[i].totalPercentage.toFixed(1));
        tableCell.appendChild(textNode);
        
        tableCell.setAttribute("title", "total percentage change over current period, base is first value from periode.");
        newRow.appendChild(tableCell);
    }
    tableContainer.appendChild(newRow);
}

function f_close_pannel() {
    let panelObj = document.getElementById("right_pannel")
    panelObj.style.display = "none";
}

function create_data_overview_right(inputIsin) {

    let panelObj = document.getElementById("right_pannel")
    panelObj.style.display = "block";

    let tableContainer = document.getElementById("overview_container");
    tableContainer.innerHTML = "";
    let currentCurveData;

    // get data that are needed
    for (let i = 0; i < gApiResponseAssetData.curves.length; i++) {
        if(gApiResponseAssetData.curves[i].name == inputIsin) {
            currentCurveData = gApiResponseAssetData.curves[i]
            break
        }        
    }

    // add data - the tilde sign means that it is editable
    // !!! important the labels with (~) are editable by the user
    // >>> the label how it is here definied is required of the mapping in the backend
    add_table_row(inputIsin, tableContainer, "Title:", currentCurveData.additionalInfo.title); 
    add_table_row(inputIsin, tableContainer, "Nickname(~):", currentCurveData.additionalInfo.nickname  );
    add_table_row(inputIsin, tableContainer, "ISIN:",  currentCurveData.name);
    add_table_row(inputIsin, tableContainer, "Color(~):", currentCurveData.additionalInfo.color  );
    add_table_row(inputIsin, tableContainer, "Country:", currentCurveData.additionalInfo.country  );
    add_table_row(inputIsin, tableContainer, "Symbol:", currentCurveData.additionalInfo.exchange + " / " + currentCurveData.additionalInfo.symbolCode);
    add_table_row(inputIsin, tableContainer, "Currency:", currentCurveData.additionalInfo.currency  );
    add_table_row(inputIsin, tableContainer, "Type:", currentCurveData.additionalInfo.type  );
    add_table_row(inputIsin, tableContainer, "Duration:", currentCurveData.additionalInfo.duration  );
    add_table_row(inputIsin, tableContainer, "Perf-1m:", currentCurveData.additionalInfo.perf1m);
    add_table_row(inputIsin, tableContainer, "Perf-6m:", currentCurveData.additionalInfo.perf6m);
    add_table_row(inputIsin, tableContainer, "Perf-1y:", currentCurveData.additionalInfo.perf1y);
    add_table_row(inputIsin, tableContainer, "Perf-5y:", currentCurveData.additionalInfo.perf5y);
    add_table_row(inputIsin, tableContainer, "Tags(~):", currentCurveData.additionalInfo.tags );
    add_table_row(inputIsin, tableContainer, "Description(~):", currentCurveData.additionalInfo.description  );
}

function add_table_row(isin, container, label, text) {

    if (label.includes("~")) {
        isWriteable = true
    } else {
        isWriteable = false
    }
    
    // Create a new row element
    const newRow = document.createElement("tr");

    // Create and append the label cell
    const labelCell = document.createElement("td");
    labelCell.appendChild(document.createTextNode(label));
    newRow.appendChild(labelCell);

    // Create and append the text cell
    const textCell = document.createElement("td");
    textCell.appendChild(document.createTextNode(text));
    
    if (isWriteable) {
        textCell.className = 'editable'
        textCell.contentEditable = 'true'
      
        // Generic event listener to update the data object on blur (editing complete)
        textCell.addEventListener('blur', (event) => {
           updateBackendAdditionalInfo(isin, label, event.target.textContent.trim()); 
        });        
    }
    
    newRow.appendChild(textCell);

    // Append the new row to the specified container
    container.appendChild(newRow);
}


async function updateBackendAdditionalInfo(isin, label, userText) {    
    console.log("under construction .... 2")
    try {
        const response = await fetch('http://localhost:8080/api/putAdditionalIsinData', {
            method: 'POST', headers: {
                'Content-Type': 'application/json'            },
                 body: JSON.stringify({     
                        isin: isin,
                        property: label,                
                        value: userText
            })        });
        // Check if response is not ok
        if (!response.ok) {
            // Attempt to extract error details from the backend response
            const errorDetails = await response.json().catch(() => ({})); // Fallback if JSON parsing fails
            throw new Error(
                `Failed to update property: ${response.status} ${response.statusText}. ` +
                `Details: ${JSON.stringify(errorDetails)}`
            );
        }
    } catch (error) {
        // Log the error for debugging and show an alert
        console.error("Error updating backend:", error);
        alert(`Error: ${error.message}`);
    }
}

function dowloadMostRecentDataFromSelectedIsins() {
    console.log("brand new function v1")

    const selectElementOptions = document.getElementById("selected-assets").options;
    const optionsArray = [...selectElementOptions];   

    for (let i = 0; i < optionsArray.length; i++) {
        let isin = optionsArray[i].innerHTML     
        downloadIsinForBackend(isin)
    }

    alert("Request update of " + optionsArray.length + " ISIN(s) send.")
}



document.addEventListener("DOMContentLoaded", function () {
    const timeModificationTable = document.getElementById("time-modification-options");
    const fromDateInput = document.getElementById("from-date");
    const toDateInput = document.getElementById("to-date");

    // Function to modify the time values
    function modifyTime(action, timeCode) {
        console.log("Start modify time function :: action_code: " + action + " / time_code: " + timeCode)

        let daysPeriod = 0
        const currentDate = new Date(); // Current date
        const resultDate = new Date(currentDate); // Copy of the current date

        const currentFromDate = new Date(fromDateInput.value);
        const currentToDate = new Date(toDateInput.value);

        // Example time modification logic based on ...
        let newFromDate = new Date(currentFromDate);
        let newToDate = new Date(currentToDate);

        switch (timeCode) {
            case '5y':
                daysPeriod = 365 * 5
                break
            case '1y':
                daysPeriod = 365
                break
            case '1m':
                daysPeriod = 31
                break
            case '7d':
                daysPeriod = 7
                break
            case '1d':
                daysPeriod = 1
                break
            default:
                console.log("action code not managed: " + action)
        }

        console.log("time periode: " + daysPeriod)

        switch (action) {
            case 'last':
                resultDate.setDate(resultDate.getDate() - daysPeriod);
                newFromDate = resultDate;
                newToDate = currentDate;
                break;
            case 'expand start':
                newFromDate.setDate(newFromDate.getDate() - daysPeriod);
                break;
            case 'shrink start':
                newFromDate.setDate(newFromDate.getDate() + daysPeriod);
                break;
            case 'scroll left':
                newFromDate.setDate(newFromDate.getDate() - daysPeriod);
                newToDate.setDate(newToDate.getDate() - daysPeriod);
                break;
            case 'scroll right':
                newFromDate.setDate(newFromDate.getDate() + daysPeriod);
                newToDate.setDate(newToDate.getDate() + daysPeriod);
                break;
            default:
                return;
        }

        // Set the new values in the input fields
        fromDateInput.value = newFromDate.toISOString().split('T')[0];
        toDateInput.value = newToDate.toISOString().split('T')[0];

        // simulate click on the Refresh button
        fetchAssetDataAndRefreshChart();
    }

    // Create buttons in every td element containing an 'x'
    const rows = timeModificationTable.rows;
    for (let i = 1; i < rows.length; i++) { // Skip the header row
        const cells = rows[i].cells;
        for (let j = 1; j < cells.length; j++) { // Skip the first column
            const cell = cells[j];
            if (cell.textContent.trim() === 'x') {
                const button = document.createElement('button');
                button.textContent = ' * ';
                button.addEventListener('click', function () {
                    const action = rows[i].cells[0].textContent.toLowerCase(); // Get the time modification type from the first column
                    const timeCode = rows[0].cells[j].textContent.toLowerCase(); // Get the time modification type from the first column
                    modifyTime(action, timeCode);
                });
                cell.textContent = ''; // Clear the 'x'
                cell.appendChild(button); // Add the button
            }
        }
    }
});
