import store from '../store.js';
import config from '../config.js';

export default class CMSFile {

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

		try {
			return fetch(`${config.api}/directories/${directory}/files`, {mode: "cors"})
				.then((response) => {

					if (response.status !== 200) {
						console.error('Oops, there was a problem', response.status);
					}

					response.json().then((data) => {
						// map documents
						store.state.documents = data.map((d) => {
							return new CMSFile(d.author, d.email, d.path, d.filename, d.title, null);
						});
					});

				});
		}
		catch(err) {
			console.error(`There was a problem retrieving documents from ${directory}, ${err}`);
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

		try {
			return fetch(path, {mode: "cors"})
				.then((response) => {

					if (response.status !== 200) {
						console.error('Oops, there was a problem', response.status);
					}

					response.json().then((d) => {
						store.state.activeDocument = new CMSFile(
							d.author, d.email, d.path, d.filename, d.title, d.html, d.markdown
						);
					});
				});
		}
		catch(err) {
			console.error(`There was a problem retrieving document from ${filename} from ${directory}, ${err}`);
		}
	};

	// instance methods

	create(commitMessage) {
		// create a commit object containing relevant info
		// and despatch it
		console.debug("creating...");
	};

	update(commit) {
		// create a commit object containing relevant info
		// and despatch it

		if (!this.changed) {
			console.warn("Update called but content hasn't changed");
		}

		var path = `${config.api}/directories/${this.path}/files/${this.filename}`

		try {

			return fetch(path, {mode: "cors", method: "PATCH", body: commit.toJSON(this)})

				.then((response) => {

					if (response.status !== 200) {
						console.error('Oops, there was a problem', response.status);
					}

					console.debug(response);


				});

		}
		catch(err) {
			console.error(`There was a problem retrieving document from ${filename} from ${directory}, ${err}`);
		}
	};

	// has the markdown changed since loading?
	get changed() {
		return this.markdown !== this.initialMarkdown;
	}

	// path/filename.md
	get absolutePath() {
		return [this.path, this.filename].join("/");
	};

}
