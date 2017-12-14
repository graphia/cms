import store from '../store.js';

export default class CMSCommit {

	constructor(message) {
		this.message = message;
	};

	static initialize() {
		store.state.commit = new CMSCommit(null);
	};

	filesJSON(file, includeAttachments=true) {

		return JSON.stringify({
			message: this.message,
			repository_info: {
				latest_revision: store.state.latestRevision
			},
			files: this._buildFilesArray(file, includeAttachments)
		});
	};

	directoriesJSON(directory) {
		return JSON.stringify({
			//message: "creating dir",
			repository_info: {
				latest_revision: store.state.latestRevision
			},
			directories: this._buildDirectoriesArray(directory)
		})
	};

	_buildFilesArray(file, includeAttachments) {
		let arr = [this._file(file)];

		if (includeAttachments) {
			arr.concat(this._attachments(file));
		};

		return arr;
	}

	_buildDirectoriesArray(directory) {
		// FIXME (maybe), only works for one directory
		return [this._directory(directory)];
	}

	_file(file) {

		let json = {
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

		return json;
	};

	_directory(directory) {
		return {
			name: directory.path,
			info: {
				title: directory.title,
				description: directory.description,
				body: directory.body
			}
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