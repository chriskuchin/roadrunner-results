import { setAuthHeader } from "../auth"

async function recordFinish(raceId, eventId, finish) {
  const url = `/api/v1/races/${raceId}/events/${eventId}/results`
  const res = await fetch(
    url, await setAuthHeader({
      method: "POST",
      headers: {
        "Content-Type": "application/json"
      },
      body: JSON.stringify({
        end_ts: finish
      })
    }))

  if (!res.ok) {
    throw new Error(`failed to record finisher: ${res.status}`)
  }
}

async function deleteResult(raceId, eventId, resultId) {
  const url = `/api/v1/races/${raceId}/events/${eventId}/results/${resultId}`
  const res = await fetch(url, await setAuthHeader({
    method: "DELETE"
  }))

  if (!res.ok) {
    throw new Error(`failed to delete result: ${res.status}`)
  }
}

export { recordFinish, deleteResult }