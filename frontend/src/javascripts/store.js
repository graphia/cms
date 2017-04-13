import Vue from 'vue';
import Vuex from 'vuex';
import router from './app';

import CMSFile from '../javascripts/models/file.js';

Vue.use(Vuex);

const state = {
	documents: []
};
const mutations = {};
const getters = {};
const actions = {
	getDocumentsInDirectory(context, directory) {
		CMSFile.all(directory);
	}
};

export default new Vuex.Store({
	actions,
	getters,
	mutations,
	state
});
