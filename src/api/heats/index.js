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
  const attemptResponse = await fetch(
    `/api/v1/races/${raceID}/events/${eventID}/heats`,
    await setAuthHeader({
      method: "POST",
      body: JSON.stringify({
        assignments: assignments,
      }),
    }),
  );

  if (attemptResponse.ok) return await attemptResponse.json()

  return []
}

export { listHeats, createNewHeat }