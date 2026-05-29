export async function listItems() {
	const response = await fetch("/api/items");
	if (!response.ok) throw new Error(await errorText(response));
	return response.json();
}

export async function createItem(input) {
	const response = await fetch("/api/items", {
		method: "POST",
		headers: { "Content-Type": "application/json" },
		body: JSON.stringify(input),
	});
	if (!response.ok) throw new Error(await errorText(response));
	return response.json();
}

export async function updateStatus(id, status) {
	const response = await fetch(`/api/items/${id}/status`, {
		method: "POST",
		headers: { "Content-Type": "application/json" },
		body: JSON.stringify({ status }),
	});
	if (!response.ok) throw new Error(await errorText(response));
	return response.json();
}

async function errorText(response) {
	try {
		const data = await response.json();
		return data.error || response.statusText;
	} catch {
		return response.statusText;
	}
}
