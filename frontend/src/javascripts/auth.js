import config from './config.js';
import {router} from './app.js';

var jwtDecode = require('jwt-decode');

export default class CMSAuth {

	constructor() {
		this._token = localStorage.getItem("token");
		console.debug("auth object created");
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
		return (this.expiry < Date.now);
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

		// is the destination somewhere other than the login page?
		if (to.path !== '/cms/login') {

			// if token exists, continue, otherwise redirect to login page
			CMSAuth.isLoggedIn() ? next() :  next('/cms/login');

		}

		// destination is login page, continue
		else {
			next();
		}

	}

	isLoggedIn() {
		this.renew();
		return !this.tokenExpired()
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
			console.debug("No token found, rendering login");
			router.push({name: 'login'});

		}
	}

};