document.addEventListener('DOMContentLoaded', function () {
    const source = new EventSource("/chart-data");
    const ctx = document.getElementById('chart').getContext('2d');

    let chart = CreateChart(ctx)

    source.onmessage = function (event) {
        let data = JSON.parse(event.data);

        // Check if the dataset for the source already exists
        chart.updateChart(data)
    };
});