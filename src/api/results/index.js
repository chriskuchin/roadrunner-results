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

async function recordResult(raceId, eventId, bib, timerId) {
  const payload = {
    bib_number: bib,
  };

  if (timerId !== "" && timerId !== "latest") {
    payload.timer_id = timerId;
  }

  const url = `/api/v1/races/${raceId}/events/${eventId}/results`;

  const res = await fetch(
    url,
    await setAuthHeader({
      method: "PUT",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(payload),
    }),
  );

  return res.ok;
}

async function updateResult(raceId, eventId, resultId, result, bib) {
  const payload = {
    result: result,
    bib_number: bib,
  }

  const url = `/api/v1/races/${raceId}/events/${eventId}/results/${resultId}`;
  const res = await fetch(
    url,
    await setAuthHeader({
      method: "PATCH",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(payload),
    }),
  );
}

async function getHeatResults(raceId, eventId, timerId) {
  const url = `/api/v1/races/${raceId}/events/${eventId}/results?timerId=${timerId}`
  const res =  await fetch(url, await setAuthHeader({
    method: "GET"
  }))

  if (!res.ok)
    return []

  return await res.json()

}


export { recordFinish, deleteResult, recordResult, updateResult, getHeatResults }