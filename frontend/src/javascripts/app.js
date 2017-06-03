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
	// {path: '/cms/signup', component: Signup, name: 'signup'},


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
router.beforeEach(store.state.auth.checkLoggedIn);

var app = new Vue({
	el: "#app",
	store,
	router,
	render: h => h(App)
});