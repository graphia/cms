import store from '../store.js';

export default class CMSFile {

	constructor(author, email, filename, title, body) {
		this.author   = author;
		this.email    = email;
		this.filename = filename;
		this.title    = title;
		this.body     = body;
	}

	// class methods

	static all(directory) {
		if (!directory){throw("directory must be specified");}

		// fetch all files in directory
		console.debug("fetching...");

		fetch(`http://localhost:8080/api/directories/${directory}/files`, {mode: "cors"})
			.then((response) => {

				if (response.status !== 200) {
					console.error('Looks like there was a problem. Status Code: ' + response.status);
				}

				response.json().then((documents) => {
					store.state.documents = documents;
				});
			});
	}

	static find(directory, filename) {
		if (!directory){throw("directory must be specified");}
		if (!filename){throw("filename must be specified");}

		console.debug(`finding ${filename} in ${directory}`);
	}

	// instance methods

	save() {
		console.debug("saving...");
	}
}
