import Vue from 'vue';
import Vuex from 'vuex';
import router from './app';

import CMSFile from '../javascripts/models/file.js';

Vue.use(Vuex);

const state = {
	documents: [],
	activeDocument: CMSFile
};
const mutations = {};
const getters = {};
const actions = {
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
