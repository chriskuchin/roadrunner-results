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

async function getParticipants(raceId, limit, offset) {
  const params = {}
  params.limit = limit || 500;
  params.offset = offset || 0;

  const res = await fetch(`/api/v1/races/${raceId}/participants?${new URLSearchParams(params).toString()}`, {
    method: "GET"
  })

  if (res.ok) {
    return await res.json()
  }

  return []
}

export { getParticipants, getParticipantByBib };
