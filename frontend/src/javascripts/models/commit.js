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

			files: [

				{
					path: document.path,
					filename: document.filename,
					body: document.markdown,

					// and the frontmatter
					frontmatter: {
						title: document.title,
						author: document.author,
						tags: document.tags,
						synopsis: document.synopsis
					}
				}

			]

		});
	};

};