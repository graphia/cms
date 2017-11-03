import store from '../store.js';
import config from '../config.js';

export default class CMSTranslation {

	constructor(path, sourceFilename, code) {
		this.path = path;
		this.sourceFilename = sourceFilename;
		this.code = code;
	};

	async create() {
		let path = `${config.api}/directories/${this.path}/files/${this.sourceFilename}/translate`

		try {
			let response = await fetch(path, {
				mode: "cors",
				method: "POST",
				headers: store.state.auth.authHeader(),
				body: this.toJSON()
			});

			return response;
		}
		catch(err) {
			console.error("There was a problem creating a translation", err);
		};
	};

	toJSON() {
		let obj = {
			source_filename: this.sourceFilename,
			path: this.path,
			language_code: this.code,
			repository_info: {
				latest_revision: store.state.latestRevision
			}
		};
		return JSON.stringify(obj);
	};


}