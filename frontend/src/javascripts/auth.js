import config from './config.js';
import {router} from './app.js';

var jwtDecode = require('jwt-decode');

export default class CMSAuth {

	constructor() {
		console.debug("Initialising CMSAuth");
		this._token = localStorage.getItem("token");
	}

	get token() {
		return this._token;
	}

	// updating the token, write the object property *and* to localStorage
	set token(value) {
		console.debug("setting token to", value);
		this._token = value;
		localStorage.setItem("token", value);
		localStorage.setItem('token_received', Date.now());
	}

	tokenExpired() {
		return (this.tokenExpiry < Date.now);
	}

	tokenExpiry() {
		if (this.token) {
			let decoded = jwtDecode(this.token);
			return decoded.exp
		}
		return 0;
	}


	// Run before every page is displayed
	checkLoggedIn(to, from, next) {

		console.debug("checking user is accessing a 'safe' path", to.path)

		// is the destination somewhere other than the login page?
		if (to.path == '/cms/login' || to.path == '/cms/setup') {
			// destination is login page, continue
			console.debug("yes they are, permit it", to.path)
			next();
		}

		else {
			console.debug("no they aren't, make sure they're logged in", to.path)

			// if token exists, continue, otherwise redirect to login page
			CMSAuth.isLoggedIn() ? next() : next('/cms/login');
		}

	}

	static async doInitialSetup() {
		let response = await fetch(`${config.auth}/create_initial_user`, {});

		console.debug("Checking for initial users!");

		if (response.status !== 200) {
			console.error('Oops, there was a problem', response.status);
			return false
		}

		let json = await response.json();

		console.debug("run initial setup:", json);

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

				console.debug("Renewing token!");

				let response = await fetch(`${config.api}/renew`,
				{mode: "cors", method: "POST", headers: this.authHeader()}
			)};
	}

	async login(username, password) {

		console.debug("logging in", username);

		let response = await fetch(`${config.auth}/login`, {
			method: "POST",
			mode: "cors",
			body: JSON.stringify({username, password})
		});

		if (response.status !== 200) {
			console.error('Oops, there was a problem', response.status);
			return false
		}

		let json = await response.json();

		// store the token and the time at which it was written
		this.token = json.token;

		return true
	}

	static async createInitialUser(user) {
		console.debug("creating initial user");

		let response = await fetch(`${config.auth}/create_initial_user`, {
			method: "POST",
			mode: "cors",
			body: JSON.stringify(user)
		});

		if (response.status !== 201) {
			console.error('Oops, there was a problem', response.status);
			return false
		}

		let json = await response.json();

		return true
		//this.redirectToLogin();
	}

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
			this.redirectToLogin();
		}
	}

	redirectToLogin() {
		// FIXME display a flash message
		router.push({name: 'login'});
	}

};