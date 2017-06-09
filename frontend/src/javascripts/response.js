import {router} from './app.js';

export default function checkResponse(responseCode) {
	console.debug("checking response", responseCode);

	if (responseCode == 401) {
		console.warn("Unauthorized request, redirecting to login");

		// Unauthorized, redirect
		router.push({name: 'login'});
		return false;

	}
	else if (responseCode < 200 || responseCode >= 300) {

		console.error('Oops, there was a problem', response.status);
		return false;

	};

	return true;

};