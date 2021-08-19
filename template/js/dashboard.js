{
    const shieldUPColor = 'rgba(0, 0, 255, 0.5)';
    const shieldDownColor = 'rgba(255, 0, 0, 0.5)';

    var users = new Map();
    var initChart = function(d, playerID, canvID, bg) {
            const data = {
            labels: [
                'Eng',
                'Weap',
                'Sys'  
            ],
        
            datasets: [{
                label: playerID,
                data: d,
                fill: true,
                backgroundColor: bg,
                borderColor: 'rgb(0, 255, 0)',
                pointBackgroundColor: 'rgb(0, 255, 0)',
                pointBorderColor: '#fff',
                pointHoverBackgroundColor: '#fff',
                pointHoverBorderColor: 'rgb(255, 99, 132)'
            }]
            };
        const config = {
            type: 'radar',
            data: data,
            options: {
                scales: {
                r: {
                    angleLines: {
                        display: true
                    },
                    suggestedMin: 0,
                    suggestedMax: 8
                }
            },
                elements: {
                line: {
                    borderWidth: 1
                }
                }
            },
            };
            
        var myChart = new Chart(
            document.getElementById(canvID),
            config
        );
        return myChart;
    };
    // var evtSource = new EventSource('https://edservertest.herokuapp.com/eventTest');
    var evtSource = new EventSource('http://localhost:1488/eventTest');


    var removeChart = function(playerID) {
        // var p = document.getElementById('main');
        var ch = document.getElementById(playerID);
        if (ch.parentNode) {
            console.log("remove", playerID);
            ch.parentNode.removeChild(ch);
        }
        ch = document.getElementById(playerID+'_status');
        if (ch.parentNode) {
            ch.parentNode.removeChild(ch);
        }
        ch = document.getElementById(playerID+'_fuel');
        if (ch.parentNode) {
            ch.parentNode.removeChild(ch);
        }
        ch = document.getElementById(playerID+'_div');
        if (ch.parentNode) {
            ch.parentNode.removeChild(ch);
        }
        users.get(playerID).destroy();
        console.log(users.has(playerID), users.delete(playerID));
        console.log(users)
        return
    }

    var addChart = function(dataEvent, x) {
        console.log("ADD")
        var d = [parseInt(dataEvent[x].Eng, 10), parseInt(dataEvent[x].Wep, 10), parseInt(dataEvent[x].Sys, 10)];
        var p = document.getElementById('main');
        var newElDiv = document.createElement('div');
        newElDiv.setAttribute('id', x+'_div');
        p.appendChild(newElDiv);

        var newEl = document.createElement('canvas');
        newEl.setAttribute('id', x);
        newElDiv.appendChild(newEl);

        var newElStatus = document.createElement('i');
        newElStatus.setAttribute('id', x+'_status');
        newElStatus.setAttribute('class', 'fas fa-balance-scale');
        newElStatus.setAttribute('style', 'display: block; padding-bottom: 2px;');
        newElStatus.innerText = ': ' + dataEvent[x].legal_state;

        var newElFuel = document.createElement('i');
        newElFuel.setAttribute('id', x+'_fuel');
        newElFuel.setAttribute('class', 'fas fa-gas-pump');

        newElFuel.innerText = ': ' + dataEvent[x].fuel_main;

        newElDiv.appendChild(newElStatus);
        newElDiv.appendChild(newElFuel);

        let bg = shieldUPColor;
        if (!dataEvent[x].shields) {
            bg = shieldDownColor;
        }
        users.set(x, initChart(d, dataEvent[x].player_name, x, bg));
    }

    var updateChart = function(dataEvent, x) {
        console.log("UPDATE")

        var d = [parseInt(dataEvent[x].Eng, 10), parseInt(dataEvent[x].Wep, 10), parseInt(dataEvent[x].Sys, 10)];
        let chart = users.get(x);
        chart.data.datasets[0].data = d;
        if (dataEvent[x].shields) {
            chart.data.datasets[0].backgroundColor = shieldUPColor;
        } else {
            chart.data.datasets[0].backgroundColor = shieldDownColor;
        }
        chart.update();
        document.getElementById(x+'_status').innerText = ': ' + dataEvent[x].legal_state;
        document.getElementById(x+'_fuel').innerText =  ': ' + dataEvent[x].fuel_main;
    }

    evtSource.onmessage = function(event) {
        let dataEvent = JSON.parse(event.data);
        if (dataEvent.delete) {
            console.log(dataEvent);
            removeChart(dataEvent.player);
            return
        }
        for (const x in dataEvent) {
            if (users.has(x)) {
                updateChart(dataEvent, x);
            } else {
                addChart(dataEvent, x);
            }
        }
    }
}