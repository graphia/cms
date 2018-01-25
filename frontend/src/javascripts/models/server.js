import config from '../config.js';
import checkResponse from '../response.js';
import store from '../store.js';

class CMSServerInfo {
	constructor() {
		this.users = null;
		this.commits = null;
		this.files = {};
	};

	async refresh() {

		let path = `${config.api}/server_info`;

		let response = await fetch(path, {
			method: "GET",
			headers: store.state.auth.authHeader()
		});

		if (!checkResponse(response.status)) {
			console.error(response);
			throw response;
		};

		let si = await response.json();

		this.commits = si.commits;
		this.users = si.users;
		this.files = si.files;

		return response;

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

class CMSRepositoryInfo {

	constructor() {
		this.latestRevision = null;
	}

	async refresh() {
		let path = `${config.api}/repository_info`;

		let response = await fetch(path, {
			method: "GET",
			headers: store.state.auth.authHeader()
		});

		if (!checkResponse(response.status)) {
			console.error(response);
		};

		let ri = await response.json();

		this.latestRevision = ri.latest_revision;

		return response;

	};

}

export default class CMSServer {

	constructor() {
		this.serverInfo = new CMSServerInfo;
		this.translationInfo = new CMSTranslationInfo;
		this.repositoryInfo = new CMSRepositoryInfo;
	};

}