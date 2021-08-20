{
    const shieldUPColor = 'rgba(0, 0, 255, 0.5)';
    const shieldDownColor = 'rgba(255, 0, 0, 0.5)';

    const activated = '100%';
    const deactivated = '5%';

    const landGearDeployed = '/ics/landing-gear-normal.png';
   
    const landedOnPlanet = '/ics/circle-normal.png';
    const dockedOnStation = '/ics/coriolis-normal.png';
    const undocked = '/ics/square-active.png';

    const supercruiseEnable = '/ics/hyperspace-normal.png';
    const fAssistEnable = '/ics/flight-assist-normal.png'

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
    var evtSource = new EventSource('https://edservertest.herokuapp.com/eventTest');
    // var evtSource = new EventSource('http://localhost:1488/eventTest');

    var removeChart = function(playerID) {
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
        var d = [parseInt(dataEvent[x].Eng, 10), parseInt(dataEvent[x].Wep, 10), parseInt(dataEvent[x].Sys, 10)];
        console.log(dataEvent[x]);
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

        var newLandGear = document.createElement('img');
        newLandGear.setAttribute('id', x+'_landing_gear');
        newLandGear.setAttribute('src', landGearDeployed);
        newLandGear.setAttribute('style', 'width: 37px; height: 35px;');
        if (dataEvent[x].land_gear) {
            newLandGear.style.opacity = activated;
        } else {
            newLandGear.style.opacity = deactivated;
        }

        var newDock = document.createElement('img');
        newDock.setAttribute('id', x+'_docked');
        if (dataEvent[x].docked) {
            newDock.src = dockedOnStation;
        } else if (dataEvent[x].landed) {
            newDock.src = landedOnPlanet;
        } else {
            newDock.src = undocked;
        }
        newDock.setAttribute('style', 'width: 37px; height: 35px;');
       

        var newSuperCruise = document.createElement('img');
        newSuperCruise.setAttribute('id', x+'_supercruise');
        newSuperCruise.setAttribute('src', supercruiseEnable);
        newSuperCruise.setAttribute('style', 'width: 37px; height: 35px;');
        if (dataEvent[x].supercruise) {
            newSuperCruise.style.opacity = activated;
        } else {
            newSuperCruise.style.opacity = deactivated;
        }

        var newFA = document.createElement('img');
        newFA.setAttribute('id', x+'_flight_assist');
        newFA.setAttribute('src', fAssistEnable);
        newFA.setAttribute('style', 'width: 37px; height: 35px;');
        if (!dataEvent[x].flight_assist) {
            newFA.style.opacity = activated;
        } else {
            newFA.style.opacity = deactivated;
        }

        newElDiv.appendChild(newLandGear);
        newElDiv.appendChild(newDock);
        newElDiv.appendChild(newSuperCruise);
        newElDiv.appendChild(newFA);

        newElDiv.appendChild(newElStatus);
        newElDiv.appendChild(newElFuel);

        let bg = shieldUPColor;
        if (!dataEvent[x].shields) {
            console.log("shield down")
            bg = shieldDownColor;
        }
        console.log(bg)

        users.set(x, initChart(d, dataEvent[x].player_name, x, bg));
    }

    var updateChart = function(dataEvent, x) {

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

        if (dataEvent[x].land_gear) {
            document.getElementById(x+'_landing_gear').style.opacity = activated;
        } else {
            document.getElementById(x+'_landing_gear').style.opacity = deactivated;
        }

        if (dataEvent[x].docked) {
            document.getElementById(x+'_docked').src = dockedOnStation;
        } else if (dataEvent[x].landed) {
            document.getElementById(x+'_docked').src = landedOnPlanet;
        } else {
            document.getElementById(x+'_docked').src = undocked;
        }

        if (dataEvent[x].supercruise) {
            document.getElementById(x+'_supercruise').style.opacity = activated;
        } else {
            document.getElementById(x+'_supercruise').style.opacity = deactivated;
        }

        if (!dataEvent[x].flight_assist) {
            document.getElementById(x+'_flight_assist').style.opacity = activated;
        } else {
            document.getElementById(x+'_flight_assist').style.opacity = deactivated;
        }
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