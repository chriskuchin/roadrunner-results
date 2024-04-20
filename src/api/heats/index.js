import { setAuthHeader } from "../auth";

async function listHeats(raceID, eventID) {
  const res = await fetch(
    `/api/v1/races/${raceID}/events/${eventID}/heats`,
    await setAuthHeader({
      method: "GET",
    }),
  );

  if (res.ok) return await res.json();

  return [];
}

async function createNewHeat(raceID, eventID, assignments) {
  const processed = []
  for (const assgnmnt of assignments) {
    processed.push({
      lane: assgnmnt.lane,
      bib: String(assgnmnt.bib)
    })
  }
  const response = await fetch(
    `/api/v1/races/${raceID}/events/${eventID}/heats`,
    await setAuthHeader({
      method: "POST",
      body: JSON.stringify({
        assignments: processed,
      }),
    }),
  );

  if (response.ok) return await response.json()

  return {}
}

async function updateHeat(raceId, eventId, heatId, assignments) {
  const processed = []
  for (const assgnmnt of assignments) {
    processed.push({
      lane: assgnmnt.lane,
      bib: String(assgnmnt.bib)
    })
  }
  const response = await fetch(
    `/api/v1/races/${raceId}/events/${eventId}/heats/${heatId}`,
    await setAuthHeader({
      method: "PUT",
      body: JSON.stringify({
        assignments: processed,
      }),
    }),
  );
}

async function deleteHeat(raceId, eventId, heatId) {
  const res = await fetch(
    `/api/v1/races/${raceId}/events/${eventId}/heats/${heatId}`,
    await setAuthHeader({
      method: "DELETE"
    })
  )
}

export { listHeats, createNewHeat, updateHeat, deleteHeat }