import Vue from 'vue';
import Vuex, { mapActions } from 'vuex';
import router from './app';

import CMSFile from '../javascripts/models/file.js';
import CMSCommit from '../javascripts/models/commit.js';
import CMSDirectory from '../javascripts/models/directory.js';
import CMSUser from '../javascripts/models/user.js';

import CMSAuth from '../javascripts/auth.js';
import CMSBroadcast from '../javascripts/broadcast.js';
import CMSServer from '../javascripts/models/server.js';

Vue.use(Vuex);

const state = {
	documents: [],
	activeDocument: new CMSFile,
	activeDirectory: new CMSDirectory,
	commit: null,
	auth: new CMSAuth,
	broadcast: new CMSBroadcast,
	user: null,
	directories: [],
	server: new CMSServer
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
		state.activeDirectory = directory;
	},
	async loadUser(context) {
		let u = await CMSUser.fetchUser();
		state.user = u;
	},
	async logout(context) {
		state.user = null;
		state.directories = []
		state.documents = []
		state.activeDocument = new CMSFile
		state.activeDirectory = new CMSDirectory
	},
	async saveUser(context, user) {
		return user.save();
	},
	async refreshServerInfo(context) {
		state.server.serverInfo.refresh();
	},
	async refreshTranslationInfo(context) {
		state.server.translationInfo.refresh();
	},
	async refreshRepositoryInfo(context) {
		state.server.repositoryInfo.refresh();
	},
	async refreshServerInfo(context) {
		state.server.serverInfo.refresh();
	},
	async setLatestRevision(context, hash) {
		state.server.repositoryInfo.latestRevision = hash;
	}
};
const getters = {};
const actions = {
	initializeCommit(context, directory) {
		return CMSCommit.initialize(directory);
	},
	getDocumentsInDirectory(context, directory) {
		state.activeDirectory = new CMSDirectory;
		return CMSFile.all(directory);
	},
	getTopLevelDirectories(context) {
		return CMSDirectory.all();
	},
	getDocument(context, args) {
		// set activeDocument to compiled doc from API
		return CMSFile.find(args.directory, args.document, args.filename);
	},
	editDocument(context, args) {
		// set activeDocument to raw doc from API
		return CMSFile.find(args.directory, args.document, args.filename, true);
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
