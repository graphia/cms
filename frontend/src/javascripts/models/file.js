import store from '../store.js';
import config from '../config.js';

export default class CMSFile {

	constructor(author, email, path, filename, title, html, markdown) {

		// TODO this is a bit ugly, can it be neatened up?
		this.author   = author;
		this.email    = email;
		this.path     = path;
		this.filename = filename;
		this.title    = title;
		this.html     = html;
		this.markdown = markdown;
	};

	// class methods

	static async all(directory) {

		try {
			return fetch(`${config.api}/directories/${directory}/files`, {mode: "cors"})
				.then((response) => {

					if (response.status !== 200) {
						console.error('Oops, there was a problem', response.status);
					}

					response.json().then((data) => {
						// map documents
						store.state.documents = data.map((d) => {
							return new CMSFile(d.author, d.email, d.path, d.filename, d.title, d.html);
						});
					});

				});
		}
		catch(err) {
			console.error(`There was a problem retrieving documents from ${directory}, ${err}`);
		}
	};

	static async find(directory, filename) {
		console.debug(`finding ${filename} in ${directory}`);

		try {
			return fetch(`${config.api}/directories/${directory}/files/${filename}`, {mode: "cors"})
				.then((response) => {

					if (response.status !== 200) {
						console.error('Oops, there was a problem', response.status);
					}

					response.json().then((d) => {
						store.state.activeDocument = new CMSFile(d.author, d.email, d.path, d.filename, d.title, d.html)
					});
				});
		}
		catch(err) {
			console.error(`There was a problem retrieving document from ${filename} from ${directory}, ${err}`);
		}

	};

	// instance methods

	save() {
		console.debug("saving...");
	};

	get absolutePath() {
		return [this.path, this.filename].join("/");
	};
}
