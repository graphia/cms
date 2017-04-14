import store from '../store.js';
import config from '../config.js';

export default class CMSFile {

	constructor(author, email, path, filename, title, body) {
		this.author   = author;
		this.email    = email;
		this.filename = filename;
		this.path     = path;
		this.title    = title;
		this.body     = body;
	}

	// class methods

	static all(directory) {
		if (!directory){throw("directory must be specified");}

		// fetch all files in directory
		console.debug("fetching...");

		fetch(`${config.api}/directories/${directory}/files`, {mode: "cors"})
			.then((response) => {

				if (response.status !== 200) {
					console.error('Looks like there was a problem. Status Code: ' + response.status);
				}

				response.json().then((data) => {
					// map documents
					//store.state.documents = documents;
					store.state.documents = data.map((d) => {
						return new CMSFile(d.author, d.email, d.filename, d.path, d.title, d.body);
					});
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

	get absolutePath() {
		return [this.path, this.filename].join("/");
	}
}
