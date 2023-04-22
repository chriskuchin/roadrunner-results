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

function formatMilliseconds(ms) {
    let min = Math.floor(ms / 60000)
    ms = ms % 60000
    let sec = Math.floor(ms / 1000)
    ms = Math.floor((ms % 1000) / 10)
    return min + ":" + addLeadingZeros(sec) + ":" + addLeadingZeros(ms)
}

function addLeadingZeros(val) {
    if (val < 10) {
        return "0" + val
    }

    return val
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
    formatMilliseconds,
    keyHandler,
}