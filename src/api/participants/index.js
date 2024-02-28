import { setAuthHeader } from "../auth";

async function getParticipantByBib(raceID, bibNumber) {
  const res = await fetch(
    `/api/v1/races/${raceID}/participants/bib/${bibNumber}`,
    await setAuthHeader({
      method: "GET",
    }),
  );

  if (res.ok) return await res.json();

  return {};
}

export { getParticipantByBib };
