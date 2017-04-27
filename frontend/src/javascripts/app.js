import Vue from 'vue';
import VueRouter from 'vue-router';

// Pages
import Login from '../components/Login.vue';
import App from '../components/App.vue';
import Home from '../components/Home.vue';
import DocumentIndex from '../components/DocumentIndex.vue';
import DocumentShow from '../components/DocumentShow.vue';
import DocumentEdit from '../components/DocumentEdit.vue';
import DocumentNew from '../components/DocumentNew.vue';


// Vuex Store
import store from './store.js';
import SimpleMDE from 'simplemde';


Vue.use(VueRouter);

const routes = [
	// Authentication
	{path: '/cms/login', component: Login},

	{path: '/cms/', component: Home},
	{path: '/cms/:directory', component: DocumentIndex, name: 'document_index'},
	{path: '/cms/:directory/new', component: DocumentNew, name: 'document_new'},
	{path: '/cms/:directory/:filename', component: DocumentShow, name: 'document_show'},
	{path: '/cms/:directory/:filename/edit', component: DocumentEdit, name: 'document_edit'}
];

const router = new VueRouter({
	routes,
	mode: 'history'
});

var app = new Vue({
	el: "#app",
	store,
	router,
	render: h => h(App)
});
