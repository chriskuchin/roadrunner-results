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


export { listTimers }