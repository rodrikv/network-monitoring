document.addEventListener('DOMContentLoaded', function () {
    const parentCanvas = document.getElementById("parentCanvas")

    fetch("/historical-data").then(function (response) {
        response.json().then(function (data) {
            data.forEach(element => {
                const ele = document.createElement("canvas");
                parentCanvas.appendChild(ele);
                const ctx = ele.getContext('2d');
                let chart = CreateChart(ctx, element.data);
                ele.ondblclick = chart.resetZoom;
            });
        });
    })
});
