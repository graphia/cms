import Vue from 'vue';
import VueRouter from 'vue-router';
import App from '../components/App.vue';
import Home from '../components/Home.vue';

import DocumentIndex from '../components/DocumentIndex.vue';

import store from './store.js';

Vue.use(VueRouter);

const routes = [
	{path: '/', component: Home},
	{path: '/:directory', component: DocumentIndex, name: 'document_index'}
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
