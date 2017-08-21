var config;

import Vue from 'vue';
import VueConfig from 'vue-config';

switch(process.env.NODE_ENV) {

	case "production":
		config = {api: "/api", auth: "/auth"};
		break;

	default:  // development
		let base_url = "http://localhost:8080";
		config = {
			api: `${base_url}/api`,
			auth: `${base_url}/auth`,
			setup: `${base_url}/setup`
		};
		break;

}

// A param named `$config` will be injected in to every component
Vue.use(VueConfig, config);

export default config;