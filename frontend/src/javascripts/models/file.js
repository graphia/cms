import store from '../store.js';
import slugify from '../utilities/slugify.js';
import config from '../config.js';
import checkResponse from '../response.js';
import CMSFileAttachment from './attachment.js';
import CMSDirectory from './directory.js';
import fecha from 'fecha';

export default class CMSFile {

	static initialize(directory) {
		let file = new CMSFile({initialzing: true, path: directory});
		store.state.activeDocument = file;
		return file;
	}

	constructor(file) {

		this.translationRegex = /\.([A-z]{2})\.md$/;

		if (file && file.initialzing) {

			this.initializing    = file.initializing;

			this.path            = file.path;
			this.document        = "";
			this.slug            = "";
			this.html            = "";
			this.markdown        = "";
			this.title           = "";
			this.author          = "";
			this.synopsis        = "";
			this.tags            = [];
			this.version         = "";
			this.history         = [];
			this.attachments     = [];
			this.translations    = [];
			this.initialMarkdown = "";
			this.language        = store.state.defaultLanguage;
			this.date            = this.todayString();
			this.draft           = true;

		} else if (file) {

			this.initializing          = false;

			// TODO this is a bit long and ugly; can it be neatened up?
			this._filename             = file.filename;
			this.path                  = file.path;
			this.document              = file.document;
			this.html                  = file.html;
			this.markdown              = file.markdown;
			this.translations          = file.translations;

			// frontmatter fields
			this.title                 = file.frontmatter.title;
			this.author                = file.frontmatter.author;
			this.synopsis              = file.frontmatter.synopsis;
			this.tags                  = file.frontmatter.tags;
			this.slug                  = file.frontmatter.slug;
			this.version               = file.frontmatter.version;
			this.draft                 = file.frontmatter.draft;
			this.date                  = file.frontmatter.date || this.todayString();

			// we don't *always* need to return directory_info with a file,
			// but if it is here, set it up
			if (file.directory_info) {
				this.directory_info = new CMSDirectory(
					this.path,
					file.directory_info.title,
					file.directory_info.description,
					file.directory_info.body
				);
			};

			// History andattachments are arrays which may be populated later
			this.history = [];     // historic commits
			this.attachments = []; // related files from the directory named after file

			// set the language by extracting the code from the filename
			this.language = this._filenameLanguage();

			// finally save the initial markdown value so we can detect changes
			// and display a diff if necessary
			this.initialMarkdown = file.markdown;

		} else {
			// do the minimum setup needed
			this.initializing   = true;
			this.draft          = true;
			this.directory_info = new CMSDirectory;
			this.translations   = [];
			this.date           = this.todayString();
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

	isTranslation() {

		if (this.language) {
			return (this.language != store.state.defaultLanguage);
		} else {
			return this.translationRegex.test(this._filename);
		};

		console.warn("Couldn't determine translation status for", this);

	};

	_filenameLanguage() {
		let code = this.translationRegex.exec(this.filename);

		if (!code) {
			return store.state.defaultLanguage;
		};

		return code[1];
	};

	set filename(value) {
		console.warn("filename cannot be set manually");
		return false;
	};

	get filename() {

		// if there is a filename value already set, use it
		if (this._filename) {
			return this._filename;
		};

		// if there's not, construct one. the format will be:
		// index.md      (for the default language)
		// index.code.md (for all other languages)
		let addLanguageCode = (store.state.translationEnabled && this.isTranslation());

		return [
			"index",
			(addLanguageCode && this.language),
			"md"
		].filter(Boolean).join(".");

	};

	get languageInfo() {
		return store.state.languages.find(x => x.code === this.language);
	};

	get attachmentsDir() {
		if (!this.path || !this.slug) {
			return null;
		};
		return [this.path, this.slug].join("/");
	};

	// Make the file usable by a commit. When deleting, we will remove
	// the entire attachments directory rather than individual attachments,
	// so includeAttachments can be set to false. In other situations, we
	// usually want to include them hence the default
	prepareJSON(includeAttachments=true) {

		let a = [];
		let f = [
			{
				path: this.path,
				filename: this.filename,
				document: this.document,
				body: this.markdown,

				// and the frontmatter
				frontmatter: {
					title: this.title,
					author: this.author,
					tags: this.tags,
					synopsis: this.synopsis,
					version: this.version,
					slug: this.slug,
					draft: this.draft,
					date: this.date
				}
			}
		];

		if (includeAttachments) {
			a = this.attachments
				.filter(attachment => attachment.isNew())
				.map((attachment) => {
					return {
						path: [this.path, this.document, "images"].join("/"),
						filename: attachment.name,
						base_64_encoded: attachment.options.base64Encoded,
						body: attachment.contents()
					};
				});
		};

		return [...f, ...a];
	};

	todayString() {
		return fecha.format(new Date, "YYYY-MM-DD");
	};

	// class methods

	/*
	 *  all retrieves all files in the specified directory and adds them to an array of CMSFile
	 *  objects, which is assigned to vuex store's `documents`
	 *  params:
	 *    directory [string]: the directory from which to retrieve files
	 */
	static async all(directory) {

		let path = `${config.api}/directories/${directory}/documents`;

		try {

			let response = await fetch(path, {mode: "cors", headers: store.state.auth.authHeader()});

			// if the api responds with a 404 we'll display a special
			// error page so handle that separately
			if (response.status == 404) {
				console.warn(`directory ${directory} not found`);
				store.state.documents = null;
				return;
			} else if (!checkResponse(response.status)) {
				// something more serious has happend, abort!
				throw(response);
			};

			let json = await response.json();

			// if we have the metadata, set up the ActiveDirectory
			if (json.info) {
				let dir = new CMSDirectory(
					directory,
					json.info.title,
					json.info.description,
					json.info.body
				);
				store.commit("setActiveDirectory", dir);
			};

			// map documents
			let docs = json.files.map((file) => {
				return new CMSFile(file);
			});

			store.state.documents = docs;
			return response;

		}
		catch(err) {
			console.error(`Couldn't retrieve files from directory ${directory}`);
		};

	};

	/*
	 *  find retrieves a single file from the API, uses its data to instantiate a CMSFile object
	 *  and assigns the object to vuex store's `activeDocument`
	 *  params:
	 *    directory [string]: the directory in which the file is stored
	 *    filename [string]: the file's filename
	 *    edit [boolean]: when true returns markdown from the API, when false HTML
	 */
	static async find(directory, document, filename, edit = false) {

		let path = `${config.api}/directories/${directory}/documents/${document}/files/${filename}`

		// if we need the uncompiled markdown (for loading the editor), amend '/edit' to the path
		if (edit) {
			path = [path, "edit"].join("/");
		};

		let response = await fetch(path, {mode: "cors", headers: store.state.auth.authHeader()})

		if (!checkResponse(response.status)) {
			console.error("Document cannot be retrieved", response);
			return;
		}

		let file = await response.json()
		let doc = new CMSFile(file);
		store.state.activeDocument = doc;

		console.debug("setting latest revision to", file.repository_info.latest_revision)
		await store.commit("setLatestRevision", file.repository_info.latest_revision);

		doc.fetchAttachments();
		return doc;

	};

	populated() {
		return (this.path || this.document);
	};

	// instance methods

	// create a commit object containing relevant info
	// and despatch it
	async create(commit) {

		if (!this.changed) {
			// TODO should probably handle this better..
			console.warn("Update called but content hasn't changed");
		}

		let path = `${config.api}/directories/${this.path}/documents`

		try {
			let response = await fetch(path, {
				mode: "cors",
				method: "POST",
				headers: store.state.auth.authHeader(),
				body: JSON.stringify(commit.prepareJSON())
			});

			return response;
		}
		catch(err) {
			console.error(`There was a problem creating new document in ${this.path}, ${err}`);
		}
	};

	async update(commit) {

		if (!this.changed) {
			console.warn("Update called but content hasn't changed");
		}

		var path = `${config.api}/directories/${this.path}/documents/${this.document}/files/${this.filename}`;

		let response = await fetch(path, {
			mode: "cors",
			method: "PATCH",
			headers: store.state.auth.authHeader(),
			body: JSON.stringify(commit.prepareJSON())
		});

		return response;

	};

	async destroy(commit) {

		var path = `${config.api}/directories/${this.path}/documents/${this.document}/files/${this.filename}`;

		let response = await fetch(path, {
			method: "DELETE",
			headers: store.state.auth.authHeader(),
			body: JSON.stringify(commit.prepareJSON(false))
		});

		return response;

	};

	async log() {

		let path = `${config.api}/directories/${this.path}/documents/${this.document}/files/${this.filename}/history`;

		let response = await fetch(path,
			{headers: store.state.auth.authHeader()}
		);

		return response;

	};

	async fetchAttachments() {

		// abort unless path and slug are present
		if (!this.path || !this.slug) {
			console.warn("Missing params, cannot retrieve attachments", this.path, this.slug);
			return;
		};

		let path = `${config.api}/directories/${this.path}/documents/${this.slug}/attachments`;

		let response = await fetch(path, {
			mode: "cors",
			method: "GET",
			headers: store.state.auth.authHeader()
		});

		if (response.status == 404) {
			console.error(`no attachments directory found for ${this.document}`);
			return;
		};

		let data = await response.json();

		if (data.length === 0) {
			console.warn(`no attachments found for ${this.document}`);
		};

		this.attachments = data.map((att) => {
			return CMSFileAttachment.fromData(att);
		});

		return;

	};

	addAttachment(attachment) {
		store.commit("addAttachment", attachment);
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
