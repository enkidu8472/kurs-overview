function executeLTRTaction(input) {
    const clipboardContentWithTRTransaction = document.getElementById("input-ltrt").value;    
    
    input.style.display = "none";  
    pushInfoToBackend(clipboardContentWithTRTransaction)
        .then(() => {
            
            input.style.display = "";               // Restore the button style after the download is complete  
            displayResponseInfo(gIsinDownloadResponse);
        })
        .catch((error) => {
            console.error("Error during TRTL action For-Backend:", error);          
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
      responseDiv.innerHTML = response
    }
  }



