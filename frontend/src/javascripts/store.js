import Vue from 'vue';
import Vuex from 'vuex';
import router from './app';

import CMSFile from '../javascripts/models/file.js';
import CMSCommit from '../javascripts/models/commit.js';
import CMSDirectory from '../javascripts/models/directory.js';

import CMSAuth from '../javascripts/auth.js';
import CMSBroadcast from '../javascripts/broadcast.js';

Vue.use(Vuex);

const state = {
	documents: [],
	activeDocument: new CMSFile,
	activeDirectory: new CMSDirectory,
	commit: null,
	auth: new CMSAuth,
	broadcast: new CMSBroadcast
};
const mutations = {
	addAttachment(context, file) {
		return state.activeDocument.attachments.push(file);
	},
	initializeDocument(context, directory) {
		let doc = CMSFile.initialize(directory);
		state.activeDocument = doc;
		return doc;
	},
	initializeDirectory(context) {
		let dir = CMSDirectory.initialize(null);
		state.activeDirectory = dir;
		return dir;
	},
	setActiveDirectory(context, directory) {
		console.debug("setting active directory", directory);
		state.activeDirectory = directory;
	}
};
const getters = {};
const actions = {
	initializeCommit(context, directory) {
		return CMSCommit.initialize(directory);
	},
	getDocumentsInDirectory(context, directory) {
		return CMSFile.all(directory);
	},
	getDocument(context, args) {
		// set activeDocument to compiled doc from API
		return CMSFile.find(args.directory, args.filename);
	},
	editDocument(context, args) {
		// set activeDocument to raw doc from API
		return CMSFile.find(args.directory, args.filename, true);
	},
	getDirectory(context, name) {
		// get and set activeDirectory by name
		return CMSDirectory.get(name);
	}
};

export default new Vuex.Store({
	actions,
	getters,
	mutations,
	state
});
