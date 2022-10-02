function formatStopwatch(min, sec, ms) {
    var minString = min;
    if (min < 10) {
        minString = "0" + min;
    }

    var secString = sec;
    if (sec < 10) {
        secString = "0" + sec;
    }

    var millisString = ms;
    if (ms < 10) {
        millisString = "0" + ms;
    }

    return minString + ":" + secString + ":" + millisString;
}

function keyHandler(charCode, handler) {
    return function (e) {
        if (e.keyCode == charCode) {
            e.preventDefault();
            handler(e);
        }
    }
}

export {
    formatStopwatch,
    keyHandler,
}