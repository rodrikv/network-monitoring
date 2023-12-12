function CreateChart(ctx) {
    let initialData = {
        labels: [],
        datasets: [],
        line: {
            tension: 0.2
        }
    };

    let chartOptions = {
        xAxes: [{
            type: 'time',
            time: {
                unit: 'second'
            }
        }],
        yAxes: [{
            ticks: {
                beginAtZero: true
            }
        }]
    }

    let chart = new Chart(ctx, {
        type: "line",
        data: initialData,
        options: chartOptions
    });

    chart.updateChart = function updateChart(data) {
        if (!Array.isArray(data)) {
            data = [data]
        }

        data.forEach(data => {
            let datasetIndex = chart.data.datasets.findIndex(dataset => dataset.label === data.source);
            if (datasetIndex === -1) {
                // If the dataset doesn't exist, create a new one
                let color = stringToColorCode(data.source)
                chart.data.datasets.push({
                    label: data.source,
                    backgroundColor: color,
                    borderColor: color,
                    tension: 0.2,
                    fill: false,
                    data: []
                });

                datasetIndex = chart.data.datasets.length - 1;
            }

            if (!chart.data.labels.includes(data.seq_id)) {
                chart.data.labels.push(data.seq_id);
            }

            // Update the data for the corresponding source
            let sourceIndex = chart.data.labels.indexOf(data.seq_id);
            chart.data.datasets[datasetIndex].data[sourceIndex] = data.response_time;

        });
        chart.update();
    }

    return chart
}