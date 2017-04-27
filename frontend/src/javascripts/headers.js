// pull the token from localStorage and return it inside a
// Headers object
function authHeader() {

	let token = localStorage.getItem('token')

	if (!token) {
		console.warn("No auth token found");
		return null;
	}

	let headers = new Headers({
		'Authorization': `Bearer ${token}`,
	});

	return headers;

};

export { authHeader };