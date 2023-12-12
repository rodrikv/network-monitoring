document.addEventListener('DOMContentLoaded', function () {
    const ctx = document.getElementById('historicalChart').getContext('2d');

    let chart = CreateChart(ctx);

    // Function to retrieve data and update the chart
    fetch("/historical-data").then(function (response) {
        response.json().then(function (data) {
            chart.updateChart(data)
        });
    })
});
