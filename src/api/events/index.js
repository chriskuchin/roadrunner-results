async function getRaceEvents(raceID) {
    let res = await fetch("/api/v1/races/" + raceID + "/events", {
        method: "GET"
    })

    if (res.ok)
        return await res.json()

    return []
}

async function deleteRaceEvent(raceID, eventID) {
    let res = await fetch("/api/v1/races/" + raceID + "/events/" + eventID, {
        method: "DELETE"
    })

}

async function createRaceEvent(raceID, description, eventType, distance) {
    console.log(raceID, description, distance, eventType)
    let res = await fetch("/api/v1/races/" + raceID + "/events", {
        method: "POST",
        body: JSON.stringify({
            description: description,
            distance: distance,
            type: eventType
        })
    })
}

export {
    getRaceEvents,
    deleteRaceEvent,
    createRaceEvent,
}