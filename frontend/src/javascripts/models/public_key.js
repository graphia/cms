import store from '../store.js';
import config from '../config.js';

export default class CMSPublicKey {

	constructor(key) {

		if (key) {
			this.id = key.id;
			this.raw = key.id;
			this.name = key.name;
			this.fingerprint = key.fingerprint;
		} else {
			this.id = -1;
			this.raw = "";
			this.name = "";
			this.fingerprint = "";
		};

	};

	static async all() {
		let path = `${config.api}/settings/ssh`;

		try {
			let response = await fetch(path, {
				method: "GET",
				headers: store.state.auth.authHeader()
			});

			return response;
		}
		catch(err) {
			console.error("There was a problem creating a new public key", err);
		};

	};

	async create() {
		let path = `${config.api}/settings/ssh`;

		try {
			let response = await fetch(path, {
				method: "POST",
				headers: store.state.auth.authHeader(),
				body: JSON.stringify({
					key: this.raw,
					name: this.name
				})
			});

			return response;
		}
		catch(err) {
			console.error("There was a problem creating a new public key", err);
		};

	};

	async delete() {
		let path = `${config.api}/settings/ssh/${this.id}`;
		try {
			let response = await fetch(path, {
				method: "DELETE",
				headers: store.state.auth.authHeader()
			});
			return response;
		}
		catch(err) {
			console.error("There was a problem deleting public key", this.id, err);
		};
	};

};