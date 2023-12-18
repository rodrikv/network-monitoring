document.addEventListener('DOMContentLoaded', function () {
    const ele = document.getElementById('historicalChart');
    const ctx = ele.getContext('2d');

    fetch("/historical-data").then(function (response) {
        response.json().then(function (data) {
            let chart = CreateChart(ctx, data);
            ele.ondblclick = chart.resetZoom;
        });
    })
});
