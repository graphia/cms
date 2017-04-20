import Vue from 'vue';
import Vuex from 'vuex';
import router from './app';

import CMSFile from '../javascripts/models/file.js';
import CMSCommit from '../javascripts/models/commit.js';

Vue.use(Vuex);

const state = {
	documents: [],
	activeDocument: new CMSFile,
	commit: new CMSCommit
};
const mutations = {};
const getters = {};
const actions = {
	initializeCommit(context, directory) {
		CMSCommit.initialize(directory)
	},
	initializeDocument(context, directory) {
		CMSFile.initialize(directory);
	},
	getDocumentsInDirectory(context, directory) {
		CMSFile.all(directory);
	},
	getDocument(context, args) {
		// returns a compiled doc
		CMSFile.find(args.directory, args.filename);
	},
	editDocument(context, args) {
		// returns a raw doc
		CMSFile.find(args.directory, args.filename, true);
	}
};

export default new Vuex.Store({
	actions,
	getters,
	mutations,
	state
});
