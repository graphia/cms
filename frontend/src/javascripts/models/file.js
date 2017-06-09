import store from '../store.js';
import config from '../config.js';
import checkResponse from '../response.js';

export default class CMSFile {

	static initialize(directory) {
		console.debug("Initialising file");
		let file = new CMSFile(null, null, directory, null, null, null, null);
		store.state.activeDocument = file;
		return file;
	}

	constructor(author, email, path, filename, title, html, markdown, synopsis, tags) {

		// TODO this is a bit long and ugly; can it be neatened up?
		this.author   = author;
		this.email    = email;
		this.path     = path;
		this.filename = filename;
		this.title    = title;
		this.html     = html;
		this.markdown = markdown;
		this.synopsis = synopsis;
		this.tags     = tags;

		// finally save the initial markdown value so we can detect changes
		// and display a diff if necessary
		this.initialMarkdown = markdown;
	};

	// class methods

	/*
	 *  all retrieves all files in the specified directory and adds them to an array of CMSFile
	 *  objects, which is assigned to vuex store's `documents`
	 *  params:
	 *    directory [string]: the directory from which to retrieve files
	 */
	static async all(directory) {

		let path = `${config.api}/directories/${directory}/files`;

		try {

			let response = await fetch(path, {mode: "cors", headers: store.state.auth.authHeader()});

			if (!checkResponse(response.status)) {
				return
			}

			let json = await response.json()

			// map documents
			store.state.documents = json.map((file) => {
				return new CMSFile(file.author, file.email, file.path, file.filename, file.title, null, null, file.synopsis, file.tags);
			});

		}
		catch(err) {
			console.error(`Couldn't retrieve files from directory ${directory}`);
		}

	};

	/*
	 *  find retrieves a single file from the API, uses its data to instantiate a CMSFile object
	 *  and assigns the object to vuex store's `activeDocument`
	 *  params:
	 *    directory [string]: the directory in which the file is stored
	 *    filename [string]: the file's filename
	 *    edit [boolean]: when true returns markdown from the API, when false HTML
	 */
	static async find(directory, filename, edit = false) {
		console.debug(`finding ${filename} in ${directory}`);

		var path = `${config.api}/directories/${directory}/files/${filename}`

		// if we need the uncompiled markdown (for loading the editor), amend '/edit' to the path
		if (edit) {
			console.debug("edit is true, adding '/edit' to", path);
			path = [path, "edit"].join("/");
		}

		let response = await fetch(path, {mode: "cors", headers: store.state.auth.authHeader()})

		 if (!checkResponse(response.status)) {
			 return
		 }

		let file = await response.json()
		store.state.activeDocument = new CMSFile(
			file.author, file.email, file.path, file.filename, file.title, file.html, file.markdown, file.synopsis, file.tags
		);

	};

	// instance methods

	// create a commit object containing relevant info
	// and despatch it
	async create(commit) {

		if (!this.changed) {
			// TODO should probably handle this better..
			console.warn("Update called but content hasn't changed");
		}

		var path = `${config.api}/directories/${this.path}/files`

		try {
			let response = await fetch(path, {
				mode: "cors",
				method: "POST",
				headers: store.state.auth.authHeader(),
				body: commit.toJSON(this)
			});

			if (!checkResponse(response.status)) {
				return
			}

			return response;
		}
		catch(err) {
			console.error(`There was a problem creating new document in ${this.path}, ${err}`);
		}
	};

	async update(commit) {
		// create a commit object containing relevant info
		// and despatch it

		if (!this.changed) {
			console.warn("Update called but content hasn't changed");
		}

		var path = `${config.api}/directories/${this.path}/files/${this.filename}`

		try {
			let response = await fetch(path, {
				mode: "cors",
				method: "PATCH",
				headers: store.state.auth.authHeader(),
				body: commit.toJSON(this)
			});

			if (!checkResponse(response.status)) {
				return
			}

			return response;
		}
		catch(err) {
			console.error(`There was a problem updating document ${filename} in ${directory}, ${err}`);
		}
	};

	destroy(commit) {
		console.debug(commit);

		var path = `${config.api}/directories/${this.path}/files/${this.filename}`

		try {
			return fetch(path, {mode: "cors", method: "DELETE", headers: store.state.auth.authHeader(), body: commit.toJSON(this)})
				.then((response) => {
					if (!checkResponse(response.status)) {
						return
					}
				});
		}
		catch(err) {
			console.error(`There was a problem deleting document ${this.filename} from ${this.path}, ${err}`);
		}

		console.debug("Deleted")
	}

	// has the markdown changed since loading?
	get changed() {
		return this.markdown !== this.initialMarkdown;
	}

	// path/filename.md
	get absolutePath() {
		return [this.path, this.filename].join("/");
	};

}
