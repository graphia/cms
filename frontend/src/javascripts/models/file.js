import store from '../store.js';
import config from '../config.js';
import CMSAuth from '../auth.js';

export default class CMSFile {

	static initialize(directory) {
		console.debug("Initialising file");
		let file = new CMSFile(null, null, directory, null, null, null, null);
		store.state.activeDocument = file;
		return file;
	}

	constructor(author, email, path, filename, title, html, markdown) {

		// TODO this is a bit long and ugly; can it be neatened up?
		this.author   = author;
		this.email    = email;
		this.path     = path;
		this.filename = filename;
		this.title    = title;
		this.html     = html;
		this.markdown = markdown;

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

		console.log(CMSAuth.authHeader());
		let response = await fetch(path, {mode: "cors", headers: CMSAuth.authHeader()});

		if (response.status !== 200) {
			console.error('Oops, there was a problem', response.status);
		}

		let json = await response.json()

		// map documents
		store.state.documents = json.map((file) => {
			return new CMSFile(file.author, file.email, file.path, file.filename, file.title, null);
		});

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

		let response = await fetch(path, {mode: "cors", headers: CMSAuth.authHeader()})

		if (response.status !== 200) {
			console.error('Oops, there was a problem', response.status);
		}

		let file = await response.json()
		store.state.activeDocument = new CMSFile(
			file.author, file.email, file.path, file.filename, file.title, file.html, file.markdown
		);

	};

	// instance methods

	create(commit) {
		// create a commit object containing relevant info
		// and despatch it

		if (!this.changed) {
			console.warn("Update called but content hasn't changed");
		}

		var path = `${config.api}/directories/${this.path}/files`

		try {
			return fetch(path, {mode: "cors", method: "POST", headers: CMSAuth.authHeader(), body: commit.toJSON(this)})
				.then((response) => {
					if (response.status !== 200) {
						console.error('Oops, there was a problem', response.status);
					}
					return response;
				});
		}
		catch(err) {
			console.error(`There was a problem creating new document in ${directory}, ${err}`);
		}
	};

	update(commit) {
		// create a commit object containing relevant info
		// and despatch it

		if (!this.changed) {
			console.warn("Update called but content hasn't changed");
		}

		var path = `${config.api}/directories/${this.path}/files/${this.filename}`

		try {
			return fetch(path, {mode: "cors", method: "PATCH", headers: CMSAuth.authHeader(), body: commit.toJSON(this)})
				.then((response) => {
					if (response.status !== 200) {
						console.error('Oops, there was a problem', response.status);
					}
				});
		}
		catch(err) {
			console.error(`There was a problem updating document ${filename} in ${directory}, ${err}`);
		}
	};

	destroy(commit) {
		console.debug(commit);

		var path = `${config.api}/directories/${this.path}/files/${this.filename}`

		try {
			return fetch(path, {mode: "cors", method: "DELETE", headers: CMSAuth.authHeader(), body: commit.toJSON(this)})
				.then((response) => {
					if (response.status !== 200) {
						console.error('Oops, there was a problem', response.status);
					}
				});
		}
		catch(err) {
			console.error(`There was a problem deleting document ${filename} from ${directory}, ${err}`);
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
