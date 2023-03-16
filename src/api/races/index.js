async function getRaces() {
    let res = await fetch("/api/v1/races", {
        method: "GET",
    })

    if (res.ok)
        return await res.json()

    return []
}

async function createRace(description) {
    let res = await fetch("/api/v1/races", {
        method: "POST",
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify({
            name: description
        })
    })

    if (!res.ok)
        return []

    return await getRaces()
}

async function deleteRace(raceID) {
    let res = await fetch("/api/v1/races/" + raceID, {
        method: "DELETE"
    })

    if (res.ok)
        return await getRaces()

    return []
}

export {
    getRaces,
    createRace,
    deleteRace,
}