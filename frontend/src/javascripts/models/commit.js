import store from '../store.js';

export default class CMSCommit {

	constructor(message, files=[], directories=[]) {
		this.message     = message;
		this.files       = files;
		this.directories = directories;
	};

	static initialize() {
		store.state.commit = new CMSCommit(null);
	};

	addFile(file) {
		this.files.push(file);
	};

	addDirectory(dir) {
		this.directories.push(dir);
	};

	reset() {
		return this.resetFiles() && this.resetDirectories();
	};

	resetFiles() {
		this.files = [];
	};

	resetDirectories() {
		this.directories = [];
	};

	// creates JSON string for transmission, includeAttachments when true adds
	// each new (or modified) attachment as a file in the files array. When
	// deleting, set to false as the entire attchments directory is likely to be
	// removed
	prepareJSON(includeAttachments=true) {

		return {
			message: this.message,
			repository_info: {
				latest_revision: store.state.latestRevision
			},
			files: this.files
				.map((f) => {return f.prepareJSON(includeAttachments)})
				.reduce((acc,cur) => {return [...acc, ...cur]}, []),
			directories: this.directories.map((d) => {return d.prepareJSON()})
		};

	};

};