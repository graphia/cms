import store from '../store.js';
import config from '../config.js';
import checkResponse from '../response.js';

export default class CMSUser {

	constructor(data) {
		console.debug("Initialising User", data);

		if (!data) {
			this._name = "";
			this._username = "";
			this._email = "";
			this.persisted = null;

			return;
		};

		this._name = data.name;
		this._username = data.username;
		this._email = data.email;

		// non limited user fields
		this.password = undefined;

		// other stuff
		this.persisted = data;
		this.refreshInProgress = false;
	};

	// username getter/setter

	get username() {
		this._checkRefreshRequired(this._username);
		return this._username;
	};

	set username(value) {
		this._username = value;
	}

	// name getter/setter

	get name() {
		this._checkRefreshRequired(this._name);
		return this._name;
	};

	set name(value) {
		this._name = value;
	}

	// email getter/setter

	get email() {
		this._checkRefreshRequired(this._email);
		return this._email;
	};

	set email(value) {
		this._email = value;
	}

	_checkRefreshRequired(value) {
		if (value === "") {
			this.refresh();
		};
	};

	populated() {
		if (this.name != "") {
			return true;
		};

		return false;
	};

	updated() {
		if (!this.persisted) {
			return false;
		};

		return ['email', 'username', 'name']
			.some((prop) => {
				return this[prop] != this.persisted[prop];
			});

	}

	async refresh() {

		if (this.refreshInProgress) {
			console.debug("already requesting user info, skipping");
			return;
		};

		try {

			this.refreshInProgress = true;

			let response = await fetch(`${config.api}/user_info`, {
				method: "GET",
				headers: store.state.auth.authHeader()
			});

			if (!checkResponse(response.status)) {
				throw {reason: "Couldn't get user info", code: response.status}
			}

			let data = await response.json();

			this.name = data.name;
			this.username = data.username;
			this.email = data.email;
			this.persisted = data;

		}
		catch(e) {
			console.error("retrieving user information failed", e);
		}
		finally {
			this.refreshInProgress = false;
		};

	};
};