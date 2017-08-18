import store from '../store.js';
import config from '../config.js';
import checkResponse from '../response.js';
import CMSFileAttachment from './attachment.js';

export default class CMSFile {

	static initialize(directory) {
		console.debug("Initialising file in", directory);
		let file = new CMSFile({initialzing: true, path: directory});
		store.state.activeDocument = file;
		return file;
	}

	constructor(file) {

		if (file && file.initialzing) {

			this.path            = file.path;
			this.filename        = "";
			this.slug            = "";
			this.html            = "";
			this.markdown        = "";
			this.title           = "";
			this.author          = "";
			this.synopsis        = "";
			this.tags            = "";
			this.version         = "";
			this.history         = [];
			this.attachments     = [];
			this.initialMarkdown = "";

		} else if (file) {

			// TODO this is a bit long and ugly; can it be neatened up?
			this.path                  = file.path;
			this.filename              = file.filename;
			this.html                  = file.html;
			this.markdown              = file.markdown;

			// frontmatter fields
			this.title                 = file.frontmatter.title;
			this.author                = file.frontmatter.author;
			this.synopsis              = file.frontmatter.synopsis;
			this.tags                  = file.frontmatter.tags;
			this.slug                  = file.frontmatter.slug;
			this.version               = file.frontmatter.version;

			// History and attachments are arrays which may be populated later
			this.history = [];     // historic commits
			this.attachments = []; // related files from the directory named after file

			// finally save the initial markdown value so we can detect changes
			// and display a diff if necessary
			this.initialMarkdown = file.markdown;

		};

	};

	set tags(tags) {
		if (typeof tags == 'string') {
			this._tags = tags.split(",");
		} else if (tags instanceof Array) {
			this._tags = tags;
		} else if (tags) {
			console.warn("tags must be an array or a comma-separated string", tags);
		};
	}

	get tags() {
		return this._tags;
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
				return new CMSFile(file);
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
		let doc = new CMSFile(file);
		store.state.activeDocument = doc;

		doc.fetchAttachments();

	};

	populated() {
		return (this.path || this.filename);
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
			console.error(`There was a problem updating document ${this.filename} in ${this.directory}, ${err}`);
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
	};

	async log() {
		let path = `${config.api}/directories/${this.path}/files/${this.filename}/history`

		try {
			let response = await fetch(path, {
				mode: "cors",
				method: "GET",
				headers: store.state.auth.authHeader()
			});

			if (!checkResponse(response.status)) {
				return
			}

			return response;
		}
		catch(err) {
			console.error(`There was a problem retriving log for ${filename} in ${directory}, ${err}`);
		}


		console.debug("Deleted")
	};

	async fetchAttachments() {

		let path = `${config.api}/directories/${this.path}/files/${this.slug}/attachments`

		try {
			let response = await fetch(path, {
				mode: "cors",
				method: "GET",
				headers: store.state.auth.authHeader()
			});

			if (response.status == 404) {
				console.debug("No attachments found");
				return;
			}

			if (!checkResponse(response.status)) {
				return;
			};

			let data = await response.json();

			this.attachments = data.map((att) => {
				return CMSFileAttachment.fromData(att);
			});

			return;

		}
		catch(err) {
			console.error(`There was a problem retriving attachments`);
		}
	};

	addAttachment(file) {
		store.commit("addAttachment", file);
	};

	// has the markdown changed since loading?
	get changed() {
		return this.markdown !== this.initialMarkdown;
	};

	// path/filename.md
	get absolutePath() {
		return [this.path, this.filename].join("/");
	};

}
