import config from '../config.js';
import checkResponse from '../response.js';
import store from '../store.js';

export default class CMSRepo {

	constructor() {
		this.repositoryInfo = null;
	}

	static async getLatestRevision() {
		let path = `${config.api}/repository_info`;

		try {
			let response = await fetch(path, {
				mode: "cors",
				method: "GET",
				headers: store.state.auth.authHeader()
			});

			console.debug("response", response)

			if (!checkResponse(response.status)) {
				console.warn("Could not retrieve repository info", response);
				return;
			};

			let repositoryInfo = await response.json();
			console.debug("repo info retrieved", repositoryInfo);
			store.state.latestRevision = repositoryInfo.latest_revision;

			return response;

		}
		catch(err) {
			console.error("There was a problem retrieving repository information");
		};

	};

}