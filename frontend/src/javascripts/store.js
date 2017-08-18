import Vue from 'vue';
import Vuex from 'vuex';
import router from './app';

import CMSFile from '../javascripts/models/file.js';
import CMSCommit from '../javascripts/models/commit.js';

import CMSAuth from '../javascripts/auth.js';
import CMSBroadcast from '../javascripts/broadcast.js';

Vue.use(Vuex);

const state = {
	documents: [],
	activeDocument: null,
	commit: null,
	auth: new CMSAuth,
	broadcast: new CMSBroadcast
};
const mutations = {
	addAttachment(context, file) {
		return state.activeDocument.attachments.push(file);
	},
	setActiveDocument(context, document) {
		console.debug("setting active document", document);
		state.activeDocument = document;
	}
};
const getters = {};
const actions = {
	initializeCommit(context, directory) {
		return CMSCommit.initialize(directory)
	},
	initializeDocument(context, directory) {
		return CMSFile.initialize(directory);
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
	}
};

export default new Vuex.Store({
	actions,
	getters,
	mutations,
	state
});
