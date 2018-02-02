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

	prepareJSON() {

		return {
			message: this.message,
			repository_info: {
				latest_revision: store.state.server.repositoryInfo.latestRevision
			},
			files: this.files
				.map((f) => {return f.prepareJSON()})
				.reduce((acc,cur) => {return [...acc, ...cur]}, []),
			directories: this.directories.map((d) => {return d.prepareJSON()})
		};

	};

};