import config from './config.js';
import store from './store.js';
import checkResponse from './response.js';

export default class CMSPublisher {

	static async publish() {

		event.preventDefault();

		if (this.publishing) {
			console.warn("already publishing, abort");
			return;
		};

		console.debug("Attempting to publish");

		let response = await fetch(`${config.api}/publish`, {
			method: "POST",
			mode: "cors",
			headers: store.state.auth.authHeader()
		});

		if (!checkResponse(response.status)) {
			console.error("Publishing failed", response);
			return;
		};

		console.log("Publishing succeeded");

		return;

	};

};