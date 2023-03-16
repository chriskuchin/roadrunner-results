async function getRaceEvents(raceID) {
    let res = await fetch("/api/v1/races/" + raceID + "/events", {
        method: "GET"
    })

    if (res.ok)
        return await res.json()

    return []
}

async function deleteRaceEvent(raceID, eventID) {

}

async function createRaceEvent(raceID) {

}

export {
    getRaceEvents,
    deleteRaceEvent,
    createRaceEvent,
}