// TODO make these methods non-static, init when log in
// and store in system state
export default class CMSAuth {

	/*
	* Run before every page is displayed
	*/
	static checkLoggedIn(to, from, next) {

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

	static isLoggedIn() {
		// TODO kick off a ping if token is older than X hours

		return !!localStorage.getItem("token");

	}

	static async ping() {
		// renew JWT
	}

	// pull the token from localStorage and return it inside a
	// Headers object
	static authHeader() {

		let token = localStorage.getItem('token')

		if (!token) {
			console.warn("No auth token found");
			return null;
		}

		let headers = new Headers({
			'Authorization': `Bearer ${token}`,
		});

		return headers;
	}

};