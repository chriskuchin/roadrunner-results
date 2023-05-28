import { setAuthHeader } from '../auth'

async function getRaces() {
    let res = await fetch("/api/v1/races", setAuthHeader({
        method: "GET",
    }))

    if (res.ok)
        return await res.json()

    return []
}

async function createRace(description) {
    let res = await fetch("/api/v1/races", setAuthHeader({
        method: "POST",
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify({
            name: description
        })
    }))

    if (!res.ok)
        return []

    return await getRaces()
}

async function deleteRace(raceID) {
    let res = await fetch("/api/v1/races/" + raceID, setAuthHeader({
        method: "DELETE"
    }))

    if (res.ok)
        return await getRaces()
    else
        throw new Error("failed to delete race: " + res.status)

    return []
}

export {
    getRaces,
    createRace,
    deleteRace,
}