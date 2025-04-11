
function updateChart1(selectedAssets) {

    document.getElementById("chart1").innerHTML = "";

    const svg = d3.select("#chart1");
    const margin = { top: 20, right: 30, bottom: 30, left: 50 };
    const width = +svg.attr("width") - margin.left - margin.right;
    const height = +svg.attr("height") - margin.top - margin.bottom;
    const chartArea = svg.append("g").attr("transform", `translate(${margin.left},${margin.top})`);

    // Create scales
    let minY = gApiResponseAssetData.curveBoundaries.price[0]
    let maxY = gApiResponseAssetData.curveBoundaries.price[1]

    let minX = new Date(gApiResponseAssetData.curveBoundaries.date[0])
    let maxX = new Date(gApiResponseAssetData.curveBoundaries.date[1])

    const xScale = d3.scaleTime()
        .domain([minX, maxX])
        .range([0, width]);


    const yScale = d3.scaleLinear()
        .domain([minY, maxY])
        .range([height, 0]);

    // Create axes
    const xAxis = d3.axisBottom(xScale).ticks(6);
    const yAxis = d3.axisLeft(yScale).ticks(6);

    // Append axes to the chart
    chartArea.append("g")
        .attr("transform", `translate(0,${height})`)
        .call(xAxis);

    chartArea.append("g")
        .call(yAxis);


    chartArea.selectAll(".line, circle").remove(); // Clear existing lines and dots

    selectedAssets.forEach(assetName => {

        const asset = gApiResponseAssetData.curves.find(d => d.name === assetName);
        if (asset === undefined) {
            console.log("no data for curve for item: " + assetName)
            return // jump to next iteration of loop
        }

        const sanitizedClassName = asset.name.replace(/\s+/g, '-'); // Replace spaces with hyphens for valid class names

        const line = d3.line()
            .x(d => xScale(d.date))
            .y(d => yScale(d.price));

        chartArea.append("path")
            .datum(asset.values)
            .attr("class", "line")
            .attr("d", line)
            .attr("stroke", asset.color);

        chartArea.selectAll(`.dot-${sanitizedClassName}`)
            .data(asset.values)
            .enter().append("circle")
            .attr("class", `dot-${sanitizedClassName}`)
            .attr("cx", d => xScale(d.date))
            .attr("cy", d => yScale(d.price))
            .attr("r", 4)
            .attr("fill", `rgba(255, 255, 255, 0)`) // 100% transparency
            .on("mouseover", (event, d) => {
                d3.select("#tooltip")
                    .style("display", "block")
                    .style("left", (event.pageX + 10) + "px")
                    .style("top", (event.pageY - 20) + "px")
                    .html(`Date: ${d.date.toDateString()}
                        <br>Price: ${d.price} 
                        <br>ISIN: ${assetName}`);
            })
            .on("mouseout", () => {
                d3.select("#tooltip").style("display", "none");
            });
    });
}


function updateChart2(selectedAssets) {
 
    document.getElementById("chart2").innerHTML = "";

    const svg = d3.select("#chart2");
    const margin = { top: 20, right: 30, bottom: 30, left: 50 };
    const width = +svg.attr("width") - margin.left - margin.right;
    const height = +svg.attr("height") - margin.top - margin.bottom;
    const chartArea = svg.append("g").attr("transform", `translate(${margin.left},${margin.top})`);

    // Create scales
    let minY = gApiResponseAssetData.curveBoundaries.pprice[0]
    let maxY = gApiResponseAssetData.curveBoundaries.pprice[1]

    let minX = new Date(gApiResponseAssetData.curveBoundaries.date[0])
    let maxX = new Date(gApiResponseAssetData.curveBoundaries.date[1])

    const xScale = d3.scaleTime()
        .domain([minX, maxX])
        .range([0, width]);

    const yScale = d3.scaleLinear()
        .domain([minY, maxY])
        .range([height, 0]);

    // Create axes
    const xAxis = d3.axisBottom(xScale).ticks(6);
    const yAxis = d3.axisLeft(yScale).ticks(6);

    // Append axes to the chart
    chartArea.append("g")
        .attr("transform", `translate(0,${height})`)
        .call(xAxis);

    chartArea.append("g")
        .call(yAxis);


    chartArea.selectAll(".line, circle").remove(); // Clear existing lines and dots

    selectedAssets.forEach(assetName => {

        const asset = gApiResponseAssetData.curves.find(d => d.name === assetName);
        if (asset === undefined) {
            console.log("no data for curve for item2: " + assetName)
            return // jump to next iteration of loop
        }
        const sanitizedClassName = asset.name.replace(/\s+/g, '-'); // Replace spaces with hyphens for valid class names
        let assetShowName = asset.additionalInfo.nickname
        if(assetShowName == "") {
            assetShowName = asset.additionalInfo.title
        }

        const line = d3.line()
            .x(d => xScale(d.date))
            .y(d => yScale(d.percentPrice));

        chartArea.append("path")
            .datum(asset.values)
            .attr("class", "line")
            .attr("d", line)
            .attr("stroke", asset.color);

        chartArea.selectAll(`.dot-${sanitizedClassName}`)
            .data(asset.values)
            .enter().append("circle")
            .attr("class", `dot-${sanitizedClassName}`)
            .attr("cx", d => xScale(d.date))
            .attr("cy", d => yScale(d.percentPrice))
            .attr("r", 4)
            .attr("fill", `rgba(255, 255, 255, 0)`) // 100% transparency
            .on("mouseover", (event, d) => {
                console.log(d)
                d3.select("#tooltip")
                    .style("display", "block")
                    .style("left", (event.pageX + 10) + "px")
                    .style("top", (event.pageY - 20) + "px")
                    .html(`Date: ${d.date.toDateString()}
                        <br>Price: €${d.price} 
                        <br>%: ${d.percentPrice.toFixed(1)} 
                        <br>ISIN: ${assetName}
                        <br>Name: ${assetShowName}`);
            })
            .on("mouseout", () => {
                d3.select("#tooltip").style("display", "none");
            });
    });




    const points = [
        { X: new Date("2025-01-01"), Y: 5.4, amount: "1.501", color: "blue", type: "buy" },
        { X: new Date("2025-02-01"), Y: 11.6, amount: "2.300", color: "green", type: "sell" },
        { X: new Date("2025-03-01"), Y: 7.2, amount: "1.200", color: "red", type: "buy" }
    ];
    
    // Iterate over the points array and create each point
    points.forEach(pts => {
        createTransactionPoint(chartArea, xScale, yScale, pts);
    });    
}

function createTransactionPoint(chartArea, xScale, yScale, pts) {
    
    const pointX = pts.X
    const pointY = pts.Y
    const color = pts.color
    let displayText = ""
    let orientation = ""
    
    if(pts.type == "buy") {
        displayText = "-" + pts.amount
        orientation = "rotate(180)"
    } else {
        displayText = "+" + pts.amount
    }

    const triangle = d3.symbol().type(d3.symbolTriangle).size(50);   
    
    // Create a frame to show details for the transaction
    const div = d3.select("body").append("div")
        .attr("id", `div-${color}-${pointX.getTime()}`) 
        .style("position", "absolute")
        .style("display", "none")
        .style("background", "lightyellow")
        .style("padding", "3px")
        .style("border", "1px solid black")
        .style("border-radius", "5px");

    chartArea.append("path")
        .attr("d", triangle) // Triangle shape        
        .attr("transform", `translate(${xScale(pointX)}, ${yScale(pointY)}) ${orientation}`)
        .attr("fill", `${color}`)
        .style("cursor", "pointer")
        .on("click", function(event) {
            const div = d3.select(`#div-${color}-${pointX.getTime()}`);            
            const isVisible = div.style("display") === "block";
            
            div.style("display", isVisible ? "none" : "block")
                .style("left", (event.pageX + 10) + "px")
                .style("top", (event.pageY - 20) + "px")
                .style("font-size", "12px")
                .html(`${displayText}`);
        });
}
