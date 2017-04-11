import Vue from 'vue';
import VueRouter from 'vue-router';
import App from '../components/App.vue';
import Home from '../components/Home.vue';

Vue.use(VueRouter);

const routes = [
	{path: '/', component: Home}
];

const router = new VueRouter({
	routes,
	mode: 'history'
});

var app = new Vue({
	el: "#app",
	router,
	render: h => h(App)
});
