// Vue stuff
import Vue from 'vue';
import VueRouter from 'vue-router';

// Pages
import App from '../components/App.vue';

import SetupInitialUser from '../components/Setup/InitialUser.vue';
import SetupInitializeRepo from '../components/Setup/InitializeRepository.vue';

import Login from '../components/Login.vue';
import Commit from '../components/Commit.vue';
import Home from '../components/Home.vue';
import Settings from '../components/Settings.vue';
import SSHKeySettings from '../components/Settings/SSHKeySettings.vue';

// Document Paths
import DocumentIndex from '../components/Document/Index.vue';
import DocumentShow from '../components/Document/Show.vue';
import DocumentEdit from '../components/Document/Edit.vue';
import DocumentNew from '../components/Document/New.vue';
import DocumentHistory from '../components/Document/History.vue';

// Directory Paths
import DirectoryNew from '../components/Directory/New.vue';

// Utility Components
import Broadcast from '../components/Broadcast.vue';
import Octicon from '../components/Utilities/Octicon.vue';

// Authentication Helpers
import CMSAuth from './auth.js';

// Vuex Store
import store from './store.js';
import SimpleMDE from 'simplemde';
import TagsInput from 'tags-input';

// Utility libs
import vagueTime from 'vague-time';

// Vue Octicons
Vue.component('octicon', Octicon);
Vue.use(VueRouter);

const routes = [
	// Unprotected pages
	{path: '/cms/setup/initial_user', component: SetupInitialUser, name: 'initial_setup'},
	{path: '/cms/login', component: Login, name: 'login'},

	// Protected pages
	{path: '/cms/setup/initialize_repo', component: SetupInitializeRepo, name: 'initialize_repo'},
	{path: '/cms/settings', component: Settings, name: 'settings'},
	{path: '/cms/settings/keys', component: SSHKeySettings, name: 'ssh_key_settings'},

	// Directory pages
	{path: '/cms/new', component: DirectoryNew, name: 'directory_new'},

	{path: '/cms/commits/:hash', component: Commit, name: 'commit'},

	{path: '/cms/', component: Home, name: 'home'},
	{path: '/cms/:directory', component: DocumentIndex, name: 'document_index'},
	{path: '/cms/:directory/new', component: DocumentNew, name: 'document_new'},
	{path: '/cms/:directory/:filename', component: DocumentShow, name: 'document_show'},
	{path: '/cms/:directory/:filename/edit', component: DocumentEdit, name: 'document_edit'},
	{path: '/cms/:directory/:filename/history', component: DocumentHistory, name: 'document_history'}

];

const router = new VueRouter({
	routes,
	mode: 'history'
});

export {router};

// ensure that only logged-in users can continue
router.beforeEach((to, from, next) => {

	console.debug("checking user is accessing a 'safe' path", to.path)

	// is the destination somewhere other than the login page?
	if (to.path == '/cms/login' || to.path == '/cms/setup/initial_user') {
		// destination is login page, continue
		next();
	}

	else {
		console.debug("no they aren't, make sure they're logged in", to.path)

		// save original destination
		window.originalDestination = to.path;

		// if token exists, continue, otherwise redirect to login page
		if (CMSAuth.isLoggedIn()) {
			console.debug("yes, they're logged in, continue");
			next();
		} else {
			console.warn("No, redirect them to the login page");
			next(new Error("NotLoggedIn"));
		};
	}

});

router.onError((err) => {
	if (err.message == "NotLoggedIn") {
		console.debug("Not logged in, redirecting to login");
		next('/cms/login');
	}
})

Vue.filter('format_date', (value) => {
	let d = new Date(Date.parse(value));
	return d.toLocaleString();
});

Vue.filter('time_ago', (value) => {
	return vagueTime.get({
		from: Date.now(),
		to: Date.parse(value),
		units: 'ms'
	});
})

Vue.filter('capitalize', (value) => {
	try {
		return value.charAt(0).toUpperCase() + value.slice(1);
	} catch(err) {
		console.warn("cannot capitalize:", value, err);
	}
});

// Create a global Event Bus
var EventBus = new Vue()

Object.defineProperties(Vue.prototype, {
	$bus: {
		get: function () {
			return EventBus;
		}
	}
});

var app = new Vue({
	el: "#app",
	store,
	router,
	render: h => h(App)
});
