var config;

import Vue from 'vue';
import VueConfig from 'vue-config';

switch(process.env.NODE_ENV) {

	case "production":
		config = {api: "//api"};
		break;

	default:  // development
		config = {api: "http://localhost:8080/api"};
		break;

}

// A param named `$config` will be injected in to every component
Vue.use(VueConfig, config);

export default config;
