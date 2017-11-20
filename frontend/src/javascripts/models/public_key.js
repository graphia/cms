import store from '../store.js';
import config from '../config.js';

export class CMSNewPublicKey {

	constructor(raw) {
		if (raw) {
			this.raw = raw;
		} else {
			this.raw = "";
		}
	};

	async create() {
		let path = `${config.api}/settings/ssh`;

		try {
			let response = await fetch(path, {
				method: "POST",
				headers: store.state.auth.authHeader(),
				body: JSON.stringify({
					key: this.raw
				})
			});

			return response;
		}
		catch(err) {
			console.error("There was a problem creating a new public key", err);
		};

	};

};

export class CMSPublicKey {

	constructor(fingerprint, raw) {
		this.fingerprint = fingerprint;
		this.raw = raw;
	};

	static async all() {

	};

	async delete() {};

};