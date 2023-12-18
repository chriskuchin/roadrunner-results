const meters_per_mile = 1609.334;
const meters_per_kilometer = 1000;

function calculateMiles(meters) {
	return meters / meters_per_mile;
}

function calculateKilometers(meters) {
	return meters / meters_per_kilometer;
}

function calculatePerKPace(ms, meters) {
	return ms / (meters / meters_per_kilometer);
}

function calculatePerMilePace(ms, meters) {
	return ms / (meters / meters_per_mile);
}

function formatMilliseconds(ms) {
	let remainingMS = 0;
	const min = Math.floor(ms / 60000);
	remainingMS = ms % 60000;
	const sec = Math.floor(ms / 1000);
	remainingMS = Math.floor((ms % 1000) / 10);

	return `${min}:${addLeadingZeros(sec)}.${addLeadingZeros(ms)}`;
}

function addLeadingZeros(val) {
	if (val < 10) {
		return `0${val}`;
	}

	return val;
}

function keyHandler(charCode, handler) {
	return (e) => {
		if (e.keyCode === charCode) {
			e.preventDefault();
			handler(e);
		}
	};
}

export {
	calculateMiles,
	calculatePerMilePace,
	calculateKilometers,
	calculatePerKPace,
	formatMilliseconds,
	keyHandler,
};
