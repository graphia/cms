import Vue from 'vue';
import VueRouter from 'vue-router';

// Pages
import Broadcast from '../components/Broadcast.vue';
import Login from '../components/Login.vue';
import App from '../components/App.vue';
import Setup from '../components/Setup.vue';
import Home from '../components/Home.vue';
import DocumentIndex from '../components/DocumentIndex.vue';
import DocumentShow from '../components/DocumentShow.vue';
import DocumentEdit from '../components/DocumentEdit.vue';
import DocumentNew from '../components/DocumentNew.vue';

// Authentication Helpers
import CMSAuth from './auth.js';

// Vuex Store
import store from './store.js';
import SimpleMDE from 'simplemde';

Vue.use(VueRouter);

const routes = [
	// Unprotected pages
	{path: '/cms/setup', component: Setup, name: 'initial_setup'},
	{path: '/cms/login', component: Login, name: 'login'},

	// Protected pages
	{path: '/cms/', component: Home, name: 'home'},
	{path: '/cms/:directory', component: DocumentIndex, name: 'document_index'},
	{path: '/cms/:directory/new', component: DocumentNew, name: 'document_new'},
	{path: '/cms/:directory/:filename', component: DocumentShow, name: 'document_show'},
	{path: '/cms/:directory/:filename/edit', component: DocumentEdit, name: 'document_edit'}
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
	if (to.path == '/cms/login' || to.path == '/cms/setup') {
		// destination is login page, continue
		next();
	}

	else {
		console.debug("no they aren't, make sure they're logged in", to.path)

		// save original destination
		window.originalDestination = to.path;

		// if token exists, continue, otherwise redirect to login page
		CMSAuth.isLoggedIn() ? next() : next(new Error("NotLoggedIn"));
	}

});

router.onError((err) => {
	if (err.message == "NotLoggedIn") {
		console.debug("Not logged in, redirecting to login");
		next('/cms/login');
	}
})

var app = new Vue({
	el: "#app",
	store,
	router,
	render: h => h(App)
});