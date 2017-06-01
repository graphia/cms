import config from '../javascripts/config.js';
import store from '../javascripts/store.js';

export default class CMSPublisher {

	static async publish() {

		event.preventDefault();

		if (this.publishing) {
			console.warn("already publishing, abort");
			return;
		};
		console.debug("Publishing");

		let response = await fetch(`${config.api}/publish`, {
			method: "POST",
			mode: "cors",
			headers: store.state.auth.authHeader()
		});

		if (response.status != 200) {
			console.error("Failed to publish", response)
			return
		};

		console.log("Publishing succeeded");

		return;

	}


}