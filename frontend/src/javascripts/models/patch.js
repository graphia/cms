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
		let text = null;

		switch (true) {
			case this.fileUpdated():
				text = "diff-modified";
				break;
			case this.fileCreated():
				text = "diff-added";
				break;
			case this.fileDeleted():
				text = "diff-removed";
				break;
		};

		return text;
	};

};