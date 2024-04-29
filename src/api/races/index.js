import { setAuthHeader } from "../auth";

async function getRaces() {
	const res = await fetch(
		"/api/v1/races",
		await setAuthHeader({
			method: "GET",
		}),
	);

	if (res.ok) return await res.json();

	return [];
}

async function createRace(description, date) {
	const res = await fetch(
		"/api/v1/races",
		await setAuthHeader({
			method: "POST",
			headers: {
				"Content-Type": "application/json",
			},
			body: JSON.stringify({
				name: description,
				date: date,
			}),
		}),
	);

	if (!res.ok) return [];

	return await getRaces();
}

async function importRace(url, description, date) {
	const res = await fetch(
		"/api/v1/races",
		await setAuthHeader({
			method: "POST",
			headers: {
				"Content-Type": "application/json",
			},
			body: JSON.stringify({
				url: url,
				name: description,
				date: date,
			}),
		}),
	);

	if (!res.ok) return [];

	return await getRaces();
}

async function deleteRace(raceID) {
	const res = await fetch(
		`/api/v1/races/${raceID}`,
		await setAuthHeader({
			method: "DELETE",
		}),
	);

	if (res.ok) return await getRaces();

	throw new Error(`failed to delete race: ${res.status}`);
}

async function getRaceVolunteers(raceId) {
	const url = `/api/v1/races/${raceId}/volunteers`;

	console.log(url)

	const res = await fetch(url, await setAuthHeader({}));

	console.log("other")

	if (res.ok) {
		return await res.json();
	}

	return [];
}

async function getRace(raceId) {
	const url = `/api/v1/races/${raceId}`
	const res = await fetch(url, await setAuthHeader({}));

	if (res.ok) {
		return await res.json()
	}

	return {}
}

async function shareRace(raceId, emails) {
	const url = `/api/v1/races/${raceId}/volunteers`;
	const payload = {
		emails: emails,
	};

	const res = await fetch(
		url,
		await setAuthHeader({
			method: "PUT",
			body: JSON.stringify(payload),
		}),
	);

	if (res.ok) {
		// added volunteers
		// handle failed adds
	} else
		throw new Error("failed to add volunteers");
}



export { getRace, getRaces, createRace, importRace, deleteRace, getRaceVolunteers, shareRace };
