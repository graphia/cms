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

	// Convert a file retrieved from the CMS into a CMSFileAttachment
	static fromData(object) {

		console.log("extracting attachment object from", object);

		let ab = new ArrayBuffer(object.data.length);

		let ia = new Uint8Array(ab);

		for (const [i, _] of object.data) {
			ia[i] = object.data.charCodeAt(i)
		}

		// write the ArrayBuffer to a blob
		let blob = new Blob([ia], { type: object.filetype });

		// Make the blob look like a File by adding name and timestamp
		blob.lastModifiedDate = new Date();
		blob.name = object.filename;

		let attachment = new CMSFileAttachment(
			blob, `data:${object.filetype};base64,${object.data}`, {}
		);

		console.log("new obj", attachment);

		return attachment;
	}

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