import store from '../store.js';
import config from '../config.js';
import checkResponse from '../response.js';

export default class CMSTranslation {

	constructor(path, document, sourceFilename, code) {
		this.path = path;
		this.document = document;
		this.sourceFilename = sourceFilename;
		this.code = code;
	};

	async create() {
		let path = `${config.api}/directories/${this.path}/documents/${this.document}/files/${this.sourceFilename}/translate`;

		try {
			let response = await fetch(path, {
				method: "POST",
				headers: store.state.auth.authHeader(),
				body: this.prepareJSON()
			});

			if (!checkResponse(response.status)) {
				throw(response);
			};

			return response;
		}
		catch(err) {
			console.error("There was a problem creating a translation", err);
		};
	};

	prepareJSON() {
		let obj = {
			source_filename: this.sourceFilename,
			path: this.path,
			source_document: this.document,
			language_code: this.code,
			repository_info: {
				latest_revision: store.state.latestRevision
			}
		};
		return JSON.stringify(obj);
	};


}