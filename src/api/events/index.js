import { setAuthHeader } from '../auth'

async function getRaceEvents(raceID) {
    let res = await fetch("/api/v1/races/" + raceID + "/events", await setAuthHeader({
        method: "GET"
    }))

    if (res.ok)
        return await res.json()

    return []
}

async function deleteRaceEvent(raceID, eventID) {
    await fetch("/api/v1/races/" + raceID + "/events/" + eventID, await setAuthHeader({
        method: "DELETE"
    }))
}

async function createRaceEvent(raceID, description, eventType, distance) {
    await fetch("/api/v1/races/" + raceID + "/events", await setAuthHeader({
        method: "POST",
        body: JSON.stringify({
            description: description,
            distance: distance,
            type: eventType
        })
    }))
}

async function getEventResults(raceID, eventID, name, gender, team, year, timers) {
    let url = `/api/v1/races/${raceID}/events/${eventID}/results`
    let filters = new URLSearchParams()

    if (name !== "") {
        filters.append("name", name)
    }

    if (timers != null && timers.length > 0) {
        timers.forEach((timer) => filters.append("timerId", timer))
    }

    if (gender != null && gender.length > 0) {
        gender.forEach((gender) => filters.append("gender", gender))
    }

    if (team != null && team.length > 0) {
        team.forEach((team) => filters.append("team", team))
    }

    if (year != null && year.length > 0) {
        year.forEach((year) => filters.append("year", year))
    }

    let res = await fetch(url + "?" + filters.toString())

    if (!res.ok) {
        console.log("error")
        return []
    }

    return await res.json()
}


export {
    getRaceEvents,
    deleteRaceEvent,
    createRaceEvent,
    getEventResults,
}