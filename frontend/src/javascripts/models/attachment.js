export default class CMSFileAttachment {

	// Create a CMSFileAttachment object
	//
	// @param {file} file A File interface https://developer.mozilla.org/en/docs/Web/API/File
	// @param {data} data The Base64 encoded contents of the file
	constructor(file, data) {

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
	}

};