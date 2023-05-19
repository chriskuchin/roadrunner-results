import { setAuthHeader } from '../auth'

async function getRaceEvents(raceID) {
    let res = await fetch("/api/v1/races/" + raceID + "/events", setAuthHeader({
        method: "GET"
    }))

    if (res.ok)
        return await res.json()

    return []
}

async function deleteRaceEvent(raceID, eventID) {
    await fetch("/api/v1/races/" + raceID + "/events/" + eventID, setAuthHeader({
        method: "DELETE"
    }))
}

async function createRaceEvent(raceID, description, eventType, distance) {
    await fetch("/api/v1/races/" + raceID + "/events", setAuthHeader({
        method: "POST",
        body: JSON.stringify({
            description: description,
            distance: distance,
            type: eventType
        })
    }))
}

export {
    getRaceEvents,
    deleteRaceEvent,
    createRaceEvent,
}