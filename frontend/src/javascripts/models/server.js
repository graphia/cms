import config from '../config.js';
import checkResponse from '../response.js';
import store from '../store.js';

class CMSServerInfo {
	constructor() {
		this.users = null;
		this.commits = null;
		this.documents = null;
		this.attachments = null;
	};

	async refresh() {

		let path = `${config.api}/server_info`;

		try {
			// let response = await fetch(path, {
			// 	method: "GET",
			// 	headers: store.state.auth.authHeader()
			// });

			// if (!checkResponse(response.status)) {
			// 	console.error(response);
			// 	throw response;
			// };

			// let si = await response.json();


			// FIXME hardcoding the response for the moment, remove

			let si = {
				commits: 5,
				users: 3,
				documents: 20,
				attachments: 15
			};

			this.commits = si.commits;
			this.users = si.users;
			this.documents = si.documents;
			this.attachments = si.attachments;

			// return response;

		}
		catch(err) {
			console.error("There was a problem retrieving repository information", err);
		};
	};

};

class CMSTranslationInfo {
	constructor() {
		this.defaultLanguage = "en";
		this.languages = [];
		this.translationEnabled = false;
	};

	async refresh() {
		let path = `${config.api}/translation_info`;

		let response = await fetch(path, {
			method: "GET",
			headers: store.state.auth.authHeader()
		});

		if (!checkResponse(response.status)) {
			console.error(response);
		};

		let ti = await response.json();

		this.defaultLanguage = ti.default_language;
		this.languages = ti.languages;
		this.translationEnabled = ti.translation_enabled;

		return response;

	};

};

export default class CMSServer {

	constructor() {
		this.serverInfo = new CMSServerInfo;
		this.translationInfo = new CMSTranslationInfo;
	};

	static async getLatestRevision() {
		let path = `${config.api}/repository_info`;

		try {
			let response = await fetch(path, {
				method: "GET",
				headers: store.state.auth.authHeader()
			});

			if (!checkResponse(response.status)) {
				console.error(response);
				throw response;
			};

			let ri = await response.json();
			store.commit("setLatestRevision", ri.latest_revision);
			//this.repositoryInfo.latestRevision = ri.latest_revision;

			return response;

		}
		catch(err) {
			console.error("There was a problem retrieving repository information", err);
		};

	};

}