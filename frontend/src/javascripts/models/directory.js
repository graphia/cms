import store from '../store.js';
import config from '../config.js';
import checkResponse from '../response.js';

export default class CMSDirectory {

	static initialize(path) {
		let dir = new CMSDirectory(path, null, null, null);
		return dir;
	};

	constructor(path, title, description, body) {
		this.path        = path        || "";
		this.title       = title       || "";
		this.description = description || "";
		this.body        = body        || "";
	};

	async create(commit) {

		let path = `${config.api}/directories`;

		try {

			let response = await fetch(path, {
				mode: "cors",
				method: "POST",
				headers: store.state.auth.authHeader(),
				body: commit.directoriesJSON(this)
			});

			if (!checkResponse(response.status)) {
				return;
			}

			return response;

		}
		catch(err) {
			console.error("There was a problem creating the new directory", err);
		};

	};

	static async get(name) {
		let path = `${config.api}/directories/${name}`;

		try {
			let response = await fetch(path, {
				mode: "cors",
				method: "GET",
				headers: store.state.auth.authHeader()
			});

			if (!checkResponse(response.status)) {
				return;
			}

			let json = await response.json();

			console.debug("JSON RESPONSE", json);

			let dir = new CMSDirectory(dir, json.title, json.description, json.body);

			store.state.activeDirectory = dir;

			return dir;

		}
		catch(err) {
			console.error("There was a problem retrieving directory metadata", err);
		}
	};

};