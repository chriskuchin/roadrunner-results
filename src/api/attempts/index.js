import { setAuthHeader } from "../auth";

async function listEventAttempts(raceID, eventID, bibNumber) {
  const res = await fetch(
    `/api/v1/races/${raceID}/events/${eventID}/attempts/${bibNumber}`,
    await setAuthHeader({
      method: "GET",
    }),
  );

  if (res.ok) return await res.json();

  return [];
}

async function recordEventAttempt(raceID, eventID, bibNumber, attempt, distance) {
  let attemptResponse = await fetch(
    `/api/v1/races/${raceID}/events/${eventID}/attempts`,
    await setAuthHeader({
      method: "POST",
      body: JSON.stringify({
        attempt_number: attempt,
        distance: distance,
        bib: `${bibNumber}`
      }),
    }),
  );

  if (attemptResponse.ok) return await attemptResponse.json()

  return []
}

export { listEventAttempts, recordEventAttempt }