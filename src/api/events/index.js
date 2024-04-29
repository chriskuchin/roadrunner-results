import { setAuthHeader } from "../auth";

async function getRaceEvents(raceId) {
	const res = await fetch(
		`/api/v1/races/${raceId}/events`,
		await setAuthHeader({
			method: "GET",
		}),
	);

	if (res.ok) return await res.json();

	return [];
}

async function deleteRaceEvent(raceId, eventId) {
	await fetch(
		`/api/v1/races/${raceId}/events/${eventId}`,
		await setAuthHeader({
			method: "DELETE",
		}),
	);
}

async function createRaceEvent(raceId, description, eventType, distance) {
	await fetch(
		`/api/v1/races/${raceId}/events`,
		await setAuthHeader({
			method: "POST",
			body: JSON.stringify({
				description: description,
				distance: Number(distance),
				type: eventType,
			}),
		}),
	);
}

async function getEventResults(
	raceId,
	eventId,
	name,
	gender,
	team,
	year,
	timers,
	order,
) {
	const url = `/api/v1/races/${raceId}/events/${eventId}/results`;
	const filters = new URLSearchParams();

	if (name !== "") {
		filters.append("name", name);
	}

	if (order !== "") {
		filters.append("order", order)
	}

	if (timers != null && timers.length > 0) {
		for (const timer of timers) {
			filters.append("timerId", timer);
		}
	}

	if (gender != null && gender.length > 0) {
		for (const gndr of gender) {
			filters.append("gender", gndr);
		}
	}

	if (team != null && team.length > 0) {
		for (const tm of team) {
			filters.append("tm", team);
		}
	}

	if (year != null && year.length > 0) {
		for (const yr of year) {
			filters.append("year", yr);
		}
	}

	const res = await fetch(`${url}?${filters.toString()}`);

	if (!res.ok) {
		return [];
	}

	return await res.json();
}

export { getRaceEvents, deleteRaceEvent, createRaceEvent, getEventResults };
