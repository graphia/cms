import store from './store.js';

export default function checkResponse(responseCode) {

	// First deal with specific status codes

	if (responseCode == 401) {
		// Unauthorized, redirect
		console.warn("Unauthorized request, redirecting to login");
		store.state.auth.redirectToLogin();
		return false;

	}

	else if (responseCode == 409) {
		// Conflict
		console.warn("Repository out of sync, abort");
		// Don't do anything just yet, this needs to be handled in
		// the originating component
		return false;

	}

	// Second broader catch-alls for other errors

	else if (responseCode >= 500) {

				console.error('Server Error', responseCode);
				return false;

	}
	else if (responseCode >= 400) {

		console.error('Client error', responseCode);
		return false;

	};

	return true;

};