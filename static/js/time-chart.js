function CreateChart(ctx, d) {
    let domain_name;
    if (d.length != 0) {
        domain_name = d[0].source
    }

    function convertData(data) {
        return data.map(d => {
            let datetime = new Date(d.timestamp)
            return{ x: datetime.valueOf(), y: d.response_time }});
    }
    let data = {
        datasets: [{
            label: 'My First dataset',
            borderColor: stringToColorCode(domain_name),
            backgroundColor: stringToColorCode(domain_name, 1),
            pointBorderColor: stringToColorCode(domain_name),
            pointBackgroundColor: stringToColorCode(domain_name,1),
            pointBorderWidth: 1,
            data: convertData(d),
        }]
    };
    // </block:data>
    // <block:scales:2>
    let scales = {
        x: {
            position: 'bottom',
            type: 'time',
            ticks: {
                autoSkip: true,
                autoSkipPadding: 50,
                maxRotation: 0
            },
            time: {
                displayFormats: {
                    hour: 'HH:mm',
                    minute: 'HH:mm',
                    second: 'HH:mm:ss'
                }
            }
        },
        y: {
            position: 'right',
            ticks: {
                callback: (val, index, ticks) => index === 0 || index === ticks.length - 1 ? null : val,
            },
            grid: {
                borderColor: stringToColorCode(domain_name),
                color: stringToColorCode(domain_name),
            },
            title: {
                display: true,
                text: (ctx) => ctx.scale.axis + ' axis',
            }
        },
    };
    // </block:scales>

    // <block:zoom:0>
    let zoomOptions = {
        pan: {
            enabled: true,
            modifierKey: 'ctrl',
        },
        zoom: {
            drag: {
                enabled: true
            },
            mode: 'x',
        },
    };
    // </block>
    // <block:config:1>
    let config = {
        type: 'scatter',
        data: data,
        options: {
            scales: scales,
            plugins: {
                zoom: zoomOptions,
                title: {
                    display: true,
                    position: 'bottom',
                }
            },
        }
    };
    // </block:config>
    let chart = new Chart(ctx, config);
}