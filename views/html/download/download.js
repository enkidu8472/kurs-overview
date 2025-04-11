function getIsinData(input) {
    const isinValue = document.getElementById("input-isin").value;    
    const loader = document.getElementById("loader");
    loader.style.display = "block";    // show load animation
    
    input.style.display = "none";  
    downloadIsinForBackend(isinValue)
        .then(() => {
            
            input.style.display = "";               // Restore the button style after the download is complete
            loader.style.display = "none";          // remove load animation
            displayResponseInfo(gIsinDownloadResponse);
        })
        .catch((error) => {
            console.error("Error during download-Isin-For-Backend:", error);          
        });
}
    
function displayResponseInfo(response) {
    const responseDiv = document.getElementById("responseInfo");
    responseDiv.innerHTML = ""; // Clear existing content
  
    // Check if there is an error message
    if (response.ErrorMessage) {
      // Display only the error message
      const errorElement = document.createElement("p");
      errorElement.textContent = `Error: ${response.ErrorMessage}`;
      errorElement.style.color = "red";
      responseDiv.appendChild(errorElement);
    } else {
      // Display all other details
      for (const key in response) {
        if (key !== "ErrorMessage") { // Exclude error message
          const detailElement = document.createElement("p");
          detailElement.textContent = `${key}: ${response[key]}`;
          responseDiv.appendChild(detailElement);
        }
      }
    }
  }



