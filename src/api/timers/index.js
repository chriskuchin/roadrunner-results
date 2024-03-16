import { setAuthHeader } from "../auth"

async function listTimers(raceId, eventId) {
  const res = await fetch(`/api/v1/races/${raceId}/events/${eventId}/timers`, await setAuthHeader({
    method: "GET"
  }))

  if (!res.ok) {
    return []
  }

  return await res.json()
}

async function startTimer(raceId, eventId, start) {
  const url = `/api/v1/races/${raceId}/events/${eventId}/timers`

  const res = await fetch(url, await setAuthHeader({
    method: "POST",
    headers: {
      "Content-Type": "application/json"
    },
    body: JSON.stringify({
      start_ts: start
    })
  }))

  if (!res.ok)
    throw new Error(`failed to start timer: ${res.status}`)

  return (await res.json()).id
}


export { listTimers, startTimer }