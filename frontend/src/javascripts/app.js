import Vue from 'vue';
import VueRouter from 'vue-router';

// Pages
import App from '../components/App.vue';
import Home from '../components/Home.vue';
import DocumentIndex from '../components/DocumentIndex.vue';
import DocumentShow from '../components/DocumentShow.vue';
import DocumentEdit from '../components/DocumentEdit.vue';

import store from './store.js';

Vue.use(VueRouter);

const routes = [
	{path: '/', component: Home},
	{path: '/:directory', component: DocumentIndex, name: 'document_index'},
	{path: '/:directory/:filename', component: DocumentShow, name: 'document_show'},
	{path: '/:directory/:filename/edit', component: DocumentEdit, name: 'document_edit'}

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
