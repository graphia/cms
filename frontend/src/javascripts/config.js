
import Vue from 'vue';
import VueConfig from 'vue-config';

const config = {api: "/api", auth: "/auth", setup: "/setup"};

// A param named `$config` will be injected in to every component
Vue.use(VueConfig, config);

export default config;
