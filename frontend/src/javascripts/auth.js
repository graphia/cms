import store from './store.js';
import config from './config.js';
import {router} from './app.js';
import checkResponse from './response';

var jwtDecode = require('jwt-decode');

export default class CMSAuth {

	constructor() {
		this._token = localStorage.getItem("token");
	};

	get token() {
		return this._token;
	};

	// updating the token, write the object property *and* to localStorage
	set token(value) {
		this._token = value;
		localStorage.setItem('token', value);
		localStorage.setItem('token_received', Date.now());
	};

	loggedIn() {
		return (this.token && !this.tokenExpired);
	};

	tokenExpired() {
		return (this.tokenExpiry < Date.now);
	};

	tokenExpiry() {
		if (this.token) {
			let decoded = jwtDecode(this.token);
			return decoded.exp
		}
		return 0;
	};

	static async doInitialSetup() {

		// check for initial users
		let response = await fetch(`${config.setup}/create_initial_user`, {});
		if (!checkResponse(response.status)) {
			console.error('Oops, there was a problem', response.status);
			return false
		}
		let json = await response.json();
		return json.enabled;

	}

	// repetetive and ugly because 'this' isn't available due
	// to checkLoggedIn being a vue beforeEach hook and 'this' not
	// being available
	static isLoggedIn() {

		let token = localStorage.getItem('token');
		var expired = true;

		if (token) {
			let exp = jwtDecode(token).exp;
			expired = (exp > Date.now());
		}

		return (token && !expired);

	}

	async renew() {
		// renew JWT if
		if (this.token &&                                    // we have a token
			!this.tokenExpired() &&                          // that's not expired
			((Date.now - this.tokenExpiry()) < (60 * 20))) { // but expires in the next 20 mins
				// renew the token
				let path = `${config.api}/renew`;
				let response = await fetch(path, {method: "POST", headers: this.authHeader()}
			)};
	}

	async login(username, password) {

		let response = await fetch(`${config.auth}/login`, {
			method: "POST",
			body: JSON.stringify({username, password})
		});

		if (response.status !== 200) {
			console.error('Oops, there was a problem', response.status);
			return false
		}

		let json = await response.json();

		// store the token and the time at which it was written
		this.token = json.jwt.token;

		["getLatestRevision", "getTopLevelDirectories"]
			.map(func => {
				store.dispatch(func);
			});

		["refreshTranslationInfo"]
			.map(func => {
				store.commit(func);
			});

		store.commit("loadUser");


		// only pull data if we're actually logged in
		// if (this.loggedIn) {
		// 	this.fetchDirectories();
		// 	this.getRepoMetadata();
		// 	this.getTranslationInfo();
		// };

		return true
	}

	async logout() {
		let path = `${config.api}/logout`

		try {
			let response = await fetch(path, {
				method: "POST",
				headers: store.state.auth.authHeader()
			});

			if (!checkResponse(response.status)) {
				throw(response);
			};

			// clear local storage and unset the user
			localStorage.removeItem('token');
			localStorage.removeItem('token_received');

			store.commit("logout");

			this.redirectToLogin();

		}
		catch(err) {
			console.error("Couldn't log out", err)
		};
	}

	static async createInitialUser(user) {

		let path = `${config.setup}/create_initial_user`;

		try {

			let response = await fetch(path, {method: "POST",body: JSON.stringify(user)});

			if (!checkResponse(response.status)) {
				throw "request failed";
			};

			return true;
		}
		catch(error) {
			console.error('Oops, there was a problem', response.status, error);
			return false;
		};
	};

	// pull the token from localStorage and return it inside a
	// Headers object
	authHeader() {

		try{
			if (!this._token) {
				console.warn("No auth token found");
				throw("No token")
			}

			let headers = new Headers({
				'Authorization': `Bearer ${this._token}`,
			});

			return headers;

		} catch(err) {
			console.warn("No token found, rendering login", err);

			store.state.broadcast.addMessage(
				"warning",
				"You are not logged in",
				`Your session probably expired, please log in again`,
				5
			);

			this.redirectToLogin();
		};

	};

	redirectToLogin() {
		router.push({name: 'login'});
	};

};