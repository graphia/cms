import config from '../config.js';
import checkResponse from '../response.js';
import store from '../store.js';

export default class CMSServer {

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

			if (!checkResponse(response.status)) {
				console.error(response);
				throw response;
			};

			let ri = await response.json();
			store.commit("setLatestRevision", ri.latest_revision);

			return response;

		}
		catch(err) {
			console.error("There was a problem retrieving repository information", err);
		};

	};

	static async getTranslationInfo() {
		let path = `${config.api}/translation_info`;

		let response = await fetch(path, {
			mode: "cors",
			method: "GET",
			headers: store.state.auth.authHeader()
		});

		try {

			if (!checkResponse(response.status)) {
				console.error(response);
				throw response;
			};

			let translationInfo = await response.json();

			store.state.defaultLanguage = translationInfo.default_language;
			store.state.languages = translationInfo.languages;
			store.state.translationEnabled = translationInfo.translation_enabled;

			return response;
		}
		catch(err) {
			console.error("Could not retreive language information", err);
		};

		return response;
	};

}