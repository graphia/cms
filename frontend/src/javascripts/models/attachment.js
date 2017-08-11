import store from '../store.js';

export default class CMSFileAttachment {

	// Create a CMSFileAttachment object
	//
	// @param {file} file A File interface https://developer.mozilla.org/en/docs/Web/API/File
	// @param {data} data The Base64 encoded contents of the file
	constructor(file, data, options) {

		this.options = options;

		this.lastModified = file.lastModified;
		this.lastModifiedDate = file.lastModifiedDate;
		this.name = file.name;

		// TODO ensure file type and size are valid
		this.size = file.size;
		this.type = file.type;

		this.data = data;
		return this;
	};

	dataURI() {
		return this.data;
	};

	// Get rid of the base64, prefix if this attachment
	// file is encoded, otherwise return the data as is
	contents() {

		if (this.options.base64Encoded) {
			return this.data.split("base64,").pop();
		}

		return this.data;
	};

	filePath() {
		return [this.dir, "images", this.name].join('/');
	};

	relativePath() {
		return ["images", this.name].join('/');
	};

	markdownImage() {
		return `![${this.name}](${window.encodeURI(this.relativePath())})`
	};


};