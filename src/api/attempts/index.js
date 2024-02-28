import { setAuthHeader } from "../auth";

async function listEventAttempts(raceID, eventID, bibNumber) {
    const res = await fetch(
        `/api/v1/races/${raceID}/events/${eventID}/results/attempts/${bibNumber}`,
        await setAuthHeader({
            method: "GET",
        }),
    );

    if (res.ok) return await res.json();

    return [];

}


export { listEventAttempts }