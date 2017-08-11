import store from '../store.js';

export default class CMSCommit {

	constructor(message) {
		this.message = message;
	};

	static initialize() {
		console.debug("Initializing commit...");
		store.state.commit = new CMSCommit(null);
	};

	toJSON(document) {
		return JSON.stringify({
			message: this.message,
			files: [this._document(document)].concat(this._attachments(document))
		});
	};

	_document(document) {

		return {
			path: document.path,
			filename: document.filename,
			body: document.markdown,

			// and the frontmatter
			frontmatter: {
				title: document.title,
				author: document.author,
				tags: document.tags,
				synopsis: document.synopsis,
				version: document.version,
				slug: document.slug
			}
		}
	}

	_attachments(document) {
		return document.attachments.map((attachment) => {
			return {
				path: [document.path, document.slug, "images"].join("/"),
				filename: attachment.name,
				base_64_encoded: attachment.options.base64Encoded,
				body: attachment.contents()
			}
		});
	}

};