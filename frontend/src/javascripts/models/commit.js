import store from '../store.js';

export default class CMSCommit {

	constructor(message) {
		this.message = message;
	}

	static initialize() {
		console.debug("Initializing commit...");
		store.state.commit = new CMSCommit(null);
	}

	toJSON(document) {
		return JSON.stringify({
			message: this.message,
			path: document.path,
			filename: document.filename,
			body: document.markdown,

			// a few extra bits to make it work without auth

			name: "Ralph Wiggum",
			email: "ralph@springfield.org",

			// and the frontmatter
			frontmatter: {
				title: document.title,
				author: document.author
			}
		});
	}

}