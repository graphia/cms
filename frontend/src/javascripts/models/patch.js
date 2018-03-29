import Diff from 'text-diff';

export default class CMSPatch {

	constructor(hash, filename, oldFile, newFile) {
		this.hash = hash;
		this.filename = filename;
		this.oldFile = oldFile;
		this.newFile = newFile;

		this.oldFilePresent = !!this.oldFile;
		this.newFilePresent = !!this.newFile;

	};

	diff() {
		let diff = new Diff();
		let textDiff = diff.main(this.oldFile, this.newFile);
		diff.cleanupSemantic(textDiff);
		return diff.prettyHtml(textDiff);
	};

	fileUpdated() {
		return (this.oldFilePresent && this.newFilePresent);
	};

	fileCreated() {
		return (!this.oldFilePresent && this.newFilePresent);
	};

	fileDeleted() {
		return (this.oldFilePresent && !this.newFilePresent);
	};

	get icon() {

		switch (true) {
			case this.fileUpdated():
				return "diff-modified";
				break;
			case this.fileCreated():
				return "diff-added";
				break;
			case this.fileDeleted():
				return "diff-removed";
				break;
		};

		console.error("File has not been not updated, created or deleted");
	};

};