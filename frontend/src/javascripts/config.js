
import Vue from 'vue';
import VueConfig from 'vue-config';

const config = {
	api: "/api",
	admin: "/api/admin",
	auth: "/auth",
	setup: "/setup",
	cms: "/cms",
	image_extensions: [".png", ".jpg", ".jpeg", ".gif"]
};

// A param named `$config` will be injected in to every component
Vue.use(VueConfig, config);

export default config;