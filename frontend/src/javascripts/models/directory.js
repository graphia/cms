import store from '../store.js';
import config from '../config.js';
import {router} from '../app.js';
import checkResponse from '../response.js';

export default class CMSDirectory {

	static initialize(path) {
		let dir = new CMSDirectory(path, null, null, null);
		return dir;
	};

	constructor(path, title, description, body, contents = []) {
		this.path        = path        || "";
		this.title       = title       || "";
		this.description = description || "";
		this.body        = body        || "";
		this.contents    = contents;
	};

	prepareJSON() {
		return {
			name: this.path,
			info: {
				title: this.title,
				description: this.description,
				body: this.body
			}
		};
	};


	static async all() {
		let path = `${config.api}/directories`

		try {
			let response = await fetch(path, {method: "GET", headers: store.state.auth.authHeader()});

			let json = await response.json();

			if (response.status == 404 && json.message == "No repository found") {
				console.warn("No repository found, redirect to create", 404)
				// FIXME redirect to create repo
			};

			if (response.status == 400 && json.message == "Not a git repository") {
				console.warn("Directory found, not git repo", 400)
				router.push({name: 'initialize_repo'});
			};

			if (!checkResponse(response.status)) {
				console.warn("error:", response);
				return;
			};

			let dirs = json.map((dir) => {
				return new CMSDirectory(
					dir.path,
					dir.info.title,
					dir.info.description,
					dir.info.body
				);
			});

			store.state.directories = dirs;

			return dirs;

		}
		catch(err) {
			console.error("Couldn't retrieve top level directory list", err);
		};

	};

	async create(commit) {

		let path = `${config.api}/directories`;

		try {

			let response = await fetch(path, {
				method: "POST",
				headers: store.state.auth.authHeader(),
				body: JSON.stringify(commit.prepareJSON())
			});

			if (!checkResponse(response.status)) {
				console.error("could not create directory", response);
				return;
			}

			return response;

		}
		catch(err) {
			console.error("There was a problem creating the new directory", err);
		};

	};

	async destroy(commit) {

		let path = `${config.api}/directories/${this.path}`;

		let response = await fetch(path, {
			method: "DELETE",
			headers: store.state.auth.authHeader(),
			body: JSON.stringify(commit.prepareJSON())
		});

		return response;

	};

	static async get(name) {
		let path = `${config.api}/directories/${name}`;

		let response = await fetch(path, {
			method: "GET",
			headers: store.state.auth.authHeader()
		});

		if (!checkResponse(response.status)) {
			console.error("couldn't get directory info", name, response);
			return;
		}

		let json = await response.json();

		let dir = new CMSDirectory(dir, json.title, json.description, json.body);

		store.state.activeDirectory = dir;

		return dir;

	};

};