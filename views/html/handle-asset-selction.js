
document.getElementById("add-button").addEventListener("click", () => {
    console.log("click add")
    const available = document.getElementById("available-assets");
    const selected = document.getElementById("selected-assets");
    const n = available.selectedOptions.length;
    for (let i = 0; i < n; i++) {
        // move the first element to the other list, until the number of selected elements is gone 
        selected.appendChild(available.selectedOptions[0]);
    }

    fetchAssetDataAndRefreshChart();

});

document.getElementById("remove-button").addEventListener("click", () => {
    const available = document.getElementById("available-assets");
    const selected = document.getElementById("selected-assets");
    const n = selected.selectedOptions.length;
    for (let i = 0; i < n; i++) {
        // move the first element to the other list, until the number of selected elements is gone
        available.appendChild(selected.selectedOptions[0]);
    }

    fetchAssetDataAndRefreshChart();
});

document.getElementById("date-reset-button").addEventListener("click", () => {
    // Reset the date input values
    document.getElementById("from-date").value = "";
    document.getElementById("to-date").value = "";
});

document.getElementById("date-refresh-button").addEventListener("click", () => {
    fetchAssetDataAndRefreshChart();
});

function manageSelectionWindow(actionCode) {
    const box = document.getElementById('center_container');
    if(actionCode == 'init') {
        box.style.display = 'block';
        fetchAssetData4SelectWindow();
    } else if (actionCode == 'confirm') {
        updateSelection();
        fetchAssetDataAndRefreshChart();
        box.style.display = 'none';
        
    } else if (actionCode == 'cancel') {
        box.style.display = 'none';

    } else {
        console("action code in manage-selection-window not handeld:" + actionCode)
    }
}

function createSelectionWindowTable(inputData) {
    console.log("createSelectionWindow-Table 1.7")
    
    // take info from main screen
    const selectedBox = document.getElementById("selected-assets");
    let alreadyCheckedIsins = [];
    for (let i = 0; i < selectedBox.length; i++) {
        alreadyCheckedIsins.push(selectedBox[i].value);
    }

    // make table
    const tableContainer = document.getElementById("selection_window_table");
    tableContainer.innerHTML = "";
    
    output = "<tr>"
    output = output + "<th>Select</th>"
    output = output + "<th>ISIN</th>"
    output = output + "<th>Name</th>"
    output = output + "<th>p1m</th>"
    output = output + "<th>p6m</th>"
    output = output + "<th>p1y</th>"
    output = output + "<th>p5y</th>"
    output = output + "<th>Title</th>"
    output = output + "<th>Description</th>"
    output = output + "<th>Tag</th></tr>";
    for (let i = 0; i < inputData.length; i++) {

        if(alreadyCheckedIsins.includes(inputData[i].name)) {
            cbAtrribute = "checked";
        } else {
            cbAtrribute = "";
        }
        
        output = output + '<tr>';
        output = output + '<td><input type="checkbox"'+ cbAtrribute +'></td>'
        output = output + '<td>' + inputData[i].name + '</td>'
        output = output + '<td>' + (inputData[i].additionalInfo?.nickname ?? 'N/A') + '</td>'
        output = output + '<td>' + (inputData[i].additionalInfo?.perf1m ?? 'N/A')+ '</td>'
        output = output + '<td>' + (inputData[i].additionalInfo?.perf6m ?? 'N/A')+ '</td>'
        output = output + '<td>' + (inputData[i].additionalInfo?.perf1y ?? 'N/A')+ '</td>'
        output = output + '<td>' + (inputData[i].additionalInfo?.perf5y ?? 'N/A')+ '</td>'
        output = output + '<td>' + (inputData[i].additionalInfo?.title ?? 'N/A')+ '</td>'
        output = output + '<td>' + (inputData[i].additionalInfo?.description ?? 'N/A') + '</td>'
        output = output + '<td>' + (inputData[i].additionalInfo?.tags  ?? 'N/A') + '</td>'
        output = output + '</tr>'
    }
    
    tableContainer.innerHTML = output;
}

function updateSelection() {

    // get user input from the big select table
    let selectedIsins = [];
    let availableIsins = [];
    document.querySelectorAll("#selection_window_table tr").forEach(row => {
        
        if(row.cells[0].tagName === "TH") {
            // dont use the table header
            return;
        };

        let checkbox = row.querySelector("td:first-child input[type='checkbox']");
        let id = row.cells[1].textContent.trim(); // Get ISIN from the second column
        if (checkbox && checkbox.checked) {
            selectedIsins.push(id);
        } else {
            availableIsins.push(id);
        }
    });
 
    // bring info to main screen
    const availableBox = document.getElementById("available-assets");
    const selectedBox = document.getElementById("selected-assets");
    availableBox.innerHTML = "";
    selectedBox.innerHTML = "";

    // fill the 'selected Isin' box
    for (let i = 0; i < selectedIsins.length; i++) {
        let option = document.createElement("option"); 
        option.textContent = selectedIsins[i]; 
        selectedBox.appendChild(option); 
    }
    
    // fill the 'available Isin' box
    for (let i = 0; i < availableIsins.length; i++) {
        let option = document.createElement("option"); 
        option.textContent = availableIsins[i]; 
        availableBox.appendChild(option); 
    }

}

function filterSelectWindow() {
    console.log("new table filter 1112")

    tagFilterValue = document.getElementById('input_filter_tag').value.trim()

    // get user input from the big select table

    let indexTag = 0;
    document.querySelectorAll("#selection_window_table tr").forEach(row => {

        // remove hidden columns from previous iteration
        row.classList.remove('hide_class');

        // find correct index of columns
        if(row.cells[0].tagName === "TH") {
            for (let i = 0; i < row.cells.length; i++ ) {
                if( row.cells[i].innerText == 'Tag'){
                    indexTag = i
                    break
                }
            }          
        };

        if (tagFilterValue != ""){
            let tagRowValue = row.cells[indexTag].textContent.trim(); 
            if ( !tagRowValue.includes(tagFilterValue) ) {
                row.classList.add('hide_class');
            }
        }
        
        // always show table header
        if (row.cells[0].tagName === "TH") {
           row.classList.remove('hide_class');
        }
    });
}

function selectAllFromSelectWindow(checkState) {

    document.querySelectorAll("#selection_window_table tr").forEach(row => {

        if(row.cells[0].tagName != "TH") {
            if(!row.classList.contains('hide_class')) {
                let checkbox = row.querySelector("td:first-child input[type='checkbox']");
                checkbox.checked = checkState;
            }           
        }
    });
}


