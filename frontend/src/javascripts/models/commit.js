import store from '../store.js';

export default class CMSCommit {

	constructor(message) {
		this.message = message;
	};

	static initialize() {
		console.debug("Initializing commit...");
		store.state.commit = new CMSCommit(null);
	};

	filesJSON(files) {
		return JSON.stringify({
			message: this.message,
			files: this._buildFilesArray(files)
		});
	};

	directoriesJSON(directory) {
		return JSON.stringify({
			//message: "creating dir",
			directories: this._buildDirectoriesArray(directory)
		})
	};

	_buildFilesArray(file) {
		// FIXME (maybe), only works for one file + attachments
		return [
			this._file(document)
		].concat(this._attachments(file));
	}

	_buildDirectoriesArray(directory) {
		// FIXME (maybe), only works for one directory
		let da =  [this._directory(directory)];
		console.debug("dir array", da);
		return da;
	}

	_file(file) {

		return {
			path: file.path,
			filename: file.filename,
			body: file.markdown,

			// and the frontmatter
			frontmatter: {
				title: file.title,
				author: file.author,
				tags: file.tags,
				synopsis: file.synopsis,
				version: file.version,
				slug: file.slug
			}
		}
	};

	_directory(directory) {
		return {
			name: directory.path
		};
	};

	_attachments(document) {
		// FIXME filter to only get new files
		return document.attachments.map((attachment) => {
			return {
				path: [document.path, document.slug, "images"].join("/"),
				filename: attachment.name,
				base_64_encoded: attachment.options.base64Encoded,
				body: attachment.contents()
			}
		});
	};

};