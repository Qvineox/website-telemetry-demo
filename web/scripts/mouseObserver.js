const pollingInterval = 500;

let lastMouseX = -1;
let lastMouseY = -1;
let lastUpdate;
let mouseTravel = 0;

let mousePath = []

addEventListener("mousemove", (event) => {
    let x = event.pageX;
    let y = event.pageY;
    let move = 0;

    if (lastMouseX > -1) {
        move = Math.max(Math.abs(x - lastMouseX), Math.abs(y - lastMouseY));
    }

    lastMouseX = x;
    lastMouseY = y;

    if (move > 2) {
        console.debug(`moved: ${move}`)
        mousePath.push([lastMouseX, lastMouseY])
    }

    // console.debug(`moved to X: ${event.x}, Y: ${event.y}`)
});

addEventListener("click", (event) => {
    monitor.sendMousePath(mousePath)
    mousePath = []
});

setInterval(function (event) {

})

// $('html').mousemove(function(e) {
//     var mousex = e.pageX;
//     var mousey = e.pageY;
//     if (lastmousex > -1)
//         mousetravel += Math.max( Math.abs(mousex-lastmousex), Math.abs(mousey-lastmousey) );
//     lastmousex = mousex;
//     lastmousey = mousey;
// });