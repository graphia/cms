// Vue stuff
import Vue from 'vue';
import VueRouter from 'vue-router';

// Pages
import App from '../components/App.vue';

import SetupInitialUser from '../components/Setup/InitialUser.vue';
import SetupInitializeRepo from '../components/Setup/InitializeRepository.vue';
import ActivateUser from '../components/Setup/ActivateUser.vue';

import Login from '../components/Login.vue';
import Commit from '../components/Commit.vue';
import Home from '../components/Home.vue';
import MyProfile from '../components/Settings/MyProfile.vue';
import SSHKeySettings from '../components/Settings/SSHKeySettings.vue';
import ThemeSettings from '../components/Settings/Theme.vue';
import History from '../components/History.vue';

// User Paths
import UserSettings from '../components/Settings/UserSettings.vue';
import UserEdit from '../components/Settings/UserSettings/Edit.vue';
import UserNew from '../components/Settings/UserSettings/New.vue';

// Document Paths
import DocumentShow from '../components/Document/Show.vue';
import DocumentEdit from '../components/Document/Edit.vue';
import DocumentNew from '../components/Document/New.vue';
import DocumentHistory from '../components/Document/History.vue';

// Directory Paths
import DirectoryIndex from '../components/Directory/Index.vue';
import DirectoryEdit from '../components/Directory/Edit.vue';
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

Vue.directive('title', {
	inserted: (el, binding) => document.title = binding.value,
	update: (el, binding) => document.title = binding.value
});

const routes = [
	// Unprotected pages
	{path: '/cms/setup/initial_user', component: SetupInitialUser, name: 'initial_setup'},
	{path: '/cms/login', component: Login, name: 'login'},
	{path: '/cms/activate/:confirmation_key', component: ActivateUser, name: 'activate_user'},

	// Protected pages
	{path: '/cms/setup/initialize_repo', component: SetupInitializeRepo, name: 'initialize_repo'},

	{path: '/cms/settings/my_profile', component: MyProfile, name: 'my_profile', alias: '/cms/settings'},
	{path: '/cms/settings/keys', component: SSHKeySettings, name: 'ssh_key_settings'},
	{path: '/cms/settings/theme', component: ThemeSettings, name: 'theme_settings'},

	// User management
	{path: '/cms/settings/users', component: UserSettings, name: 'user_settings'},
	{path: '/cms/settings/users/new', component: UserNew, name: 'user_new'},
	{path: '/cms/settings/users/:username/edit', component: UserEdit, name: 'user_edit'},


	// Directory pages
	{path: '/cms/history', component: History, name: 'history'},
	{path: '/cms/new', component: DirectoryNew, name: 'directory_new'},

	{path: '/cms/commits/:hash', component: Commit, name: 'commit'},

	{path: '/cms/', component: Home, name: 'home'},
	{path: '/cms/:directory', component: DirectoryIndex, name: 'directory_index'},
	{path: '/cms/:directory/edit', component: DirectoryEdit, name: 'directory_edit'},
	{path: '/cms/:directory/new', component: DocumentNew, name: 'document_new'},

	{path: '/cms/:directory/:document/:language_code?/edit', component: DocumentEdit, name: 'document_edit'},
	{path: '/cms/:directory/:document/:language_code?/history', component: DocumentHistory, name: 'document_history'},
	{path: '/cms/:directory/:document/:language_code?', component: DocumentShow, name: 'document_show'},

];

const router = new VueRouter({
	routes,
	mode: 'history',
	linkActiveClass: 'active',
	scrollBehavior (to, from, savedPosition) {
		return !savedPosition ? { x: 0, y: 0 } : savedPosition
	}
});

export {router};

// ensure that only logged-in users can continue
router.beforeEach((to, from, next) => {

	let allow = CMSAuth.unblockedPageCheck(to.path);

	// is the destination somewhere other than the login page?
	if (allow) {
		// destination is login page, continue
		next();
	}

	else {
		// save original destination
		window.originalDestination = to.path;

		// if token exists, continue, otherwise redirect to login page
		if (CMSAuth.isLoggedIn()) {
			next();
		} else {
			console.warn("Redirecting to the login page");
			next(new Error("NotLoggedIn"));
		};
	}

});

router.onError((err) => {
	if (err.message == "NotLoggedIn") {
		next('/cms/login');
	}
});

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
});

Vue.filter('capitalize', (value) => {
	try {
		return value.charAt(0).toUpperCase() + value.slice(1);
	} catch(err) {
		console.warn("cannot capitalize:", value, err);
		return value;
	}
});

// Create a global Event Bus
var EventBus = new Vue();

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
