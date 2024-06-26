const meters_per_mile = 1609.334;
const meters_per_kilometer = 1000;
const centimeters_per_foot = 30.48;
const centimeters_per_inch = 2.54;

function calculateCentimeters(large, small, format) {
  switch (format) {
    case "ftin":
      return large * centimeters_per_foot + small * centimeters_per_inch
    default:
      return 0
  }
}

function formatCentimeters(cm, format) {
  switch (format) {
    case "ftin": {
      const ft = Math.floor(cm / centimeters_per_foot);
      const remainingCM = cm - (ft * centimeters_per_foot);

      const inches = remainingCM / centimeters_per_inch;

      if (ft === 0) {
        return `${inches}"`;
      }

      if (inches === 0) {
        return `${ft}'`;
      }

      return `${ft}' ${Math.round(inches * 100) / 100}"`;
    }
    case "mcm":
      return "1m 3cm";

    case "cm":
      return `${cm} cm`
    default:
      return `${cm}`;
  }
}

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

function calculateMilliseconds(formattedTime) {
  const timeArray = formattedTime.split(":");
  if (timeArray.length < 2)
    throw new Error("Invalid time format");

  const min = parseInt(timeArray[0]);

  const secondsArray = timeArray[1].split(".");

  const sec = parseInt(secondsArray[0]);
  let ms = 0
  if (secondsArray.length === 2) {
    ms = parseInt(secondsArray[1]);

    if (secondsArray[1].length === 1) {
      ms *= 100;
    } else if (secondsArray[1].length === 2) {
      ms *= 10;
    }
  }

  return (min * 60000) + (sec * 1000) + ms;
}

function formatMilliseconds(ms) {
  let remainingMS = 0;
  const min = Math.floor(ms / 60000);
  remainingMS = ms % 60000;
  const sec = Math.floor(remainingMS / 1000);
  remainingMS = Math.floor((ms % 1000) / 10);

  return `${min}:${addLeadingZeros(sec)}.${addLeadingZeros(remainingMS)}`;
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
  calculateCentimeters,
  calculateMiles,
  calculatePerMilePace,
  calculateKilometers,
  calculatePerKPace,
  formatMilliseconds,
  calculateMilliseconds,
  keyHandler,
  formatCentimeters,
};
