import store from '../store.js';
import config from '../config.js';
import checkResponse from '../response.js';

export default class CMSDirectory {

	static initialize(path) {
		let dir = new CMSDirectory(path)
		return dir;
	};

	constructor(path) {
		this.path = path || "";
	};

	async create(commit) {

		let path = `${config.api}/directories`

		try {

			let response = await fetch(path, {
				mode: "cors",
				method: "POST",
				headers: store.state.auth.authHeader(),
				body: commit.directoriesJSON(this)
			});

			if (!checkResponse(response.status)) {
				return
			}

			return response;

		}
		catch(err) {
			console.error("There was a problem creating the new directory", err);
		}

	};

};