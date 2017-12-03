<template>
	<div id="application">

		<!-- Primary Navigation Start -->
		<nav class="navbar navbar-expand-md navbar-dark bg-dark">

			<router-link class="navbar-brand" :to="{name: 'home'}">Graphia CMS</router-link>


			<button class="navbar-toggler navbar-toggler-right hidden-md-up" type="button" data-toggle="collapse" data-target="#primary" aria-label="Toggle navigation">
				<span class="navbar-toggler-icon"></span>
			</button>

			<div id="primary" class="collapse navbar-collapse">
				<ul class="navbar-nav mr-auto">

					<router-link :to="{name: 'home'}" class="nav-link home-link">Home</router-link>

					<router-link v-for="(directory, i) in directories" :key="i" :to="{name: 'document_index', params: {directory: directory.path}}" class="nav-link directory-link">
						{{ directory.path | capitalize }}
					</router-link>

					<li><a class="nav-link" href="#">History</a></li>
					<li><a class="nav-link" href="#">Admin</a></li>

				</ul>

				<ul class="navbar-nav" v-if="user">
					<li class="nav-item dropdown">
						<a class="nav-link dropdown-toggle" href="#" id="user-menu" data-toggle="dropdown" aria-haspopup="true" aria-expanded="false">
							{{ user.name }}
						</a>
						<div class="dropdown-menu" aria-labelledby="user-menu">
							<router-link :to="{name: 'user_settings'}" class="dropdown-item">
								Settings
							</router-link>
							<a class="dropdown-item logout" href="logout" @click="logout">Logout</a>
						</div>
					</li>
				</ul>
			</div>

		</nav>
		<!-- Primary Navigation End -->

		<!-- Router View Container Start -->
		<div class="container-fluid">
			<Broadcast/>
			<transition name="fade">
				<router-view/>
			</transition>
		</div>
		<!-- Router View Container End -->

	</div>
</template>

<script lang="babel">
	import store from '../javascripts/store.js';
	import config from '../javascripts/config.js';

	import checkResponse from '../javascripts/response.js';
	import CMSAuth from '../javascripts/auth.js';
	import Broadcast from '../components/Broadcast';

	export default {
		name: "GraphiaCMS",
		created() {

			try {
				if (!this.$store.state.auth.tokenExpiry()) {
					throw 'Token missing';
				}

				if (this.$store.state.auth.expiry < Date.now()) {
					throw 'Token expired';
				}

				console.debug("token is present and has not expired, renewing");
				this.$store.state.auth.renew();

				// only pull data if we're actually logged in
				if (CMSAuth.isLoggedIn()) {
					this.fetchDirectories();
					this.getRepoMetadata();
					this.getTranslationInfo();
				};

			}
			catch(err) {
				// Token rejected for renewal
				console.warn(err);
				this.$store.state.auth.redirectToLogin();
			}
		},

		data() {
			return {
				directories: []
			};
		},
		computed: {
			user() {
				if (CMSAuth.isLoggedIn() && !this.$store.state.user) {
					this.$store.commit("setUser");
				};

				return this.$store.state.user;
			},
		},
		methods: {

			redirectToInitializeRepo() {
				this.$router.push({name: 'initialize_repo'});
			},
			async getRepoMetadata() {
				this.$store.dispatch("getLatestRevision");
			},
			async getTranslationInfo() {
				this.$store.dispatch("getTranslationInfo");
			},
			async fetchDirectories() {
				let path = `${config.api}/directories`

				try {
					let response = await fetch(path, {method: "GET", mode: "cors", headers: store.state.auth.authHeader()});

					let json = await response.json();

					// FIXME is this still required?
					if (response.status == 404 && json.message == "No repository found") {
						console.warn("No repository found, redirect to create", 404)
					};

					if (response.status == 400 && json.message == "Not a git repository") {
						console.warn("Directory found, not git repo", 400)
						this.redirectToInitializeRepo();
					};

					if (!checkResponse(response.status)) {
						console.warn("error:", response);
						return;
					};

					// everything's ok, set directories to the response's json
					this.directories = json;

					return;

				}
				catch(err) {
					console.error("Couldn't retrieve top level directory list", err);
				};

			},
			logout(event) {
				event.preventDefault();
				this.$store.state.auth.logout();
			}
		},
		components: {
			Broadcast
		}
	}
</script>

<style lang="scss">
	.fade-enter-active, .fade-leave-active {
		transition-property: opacity;
		transition-duration: 0.15s;
	}

	.fade-enter-active {
		transition-delay: 0.15s;
	}

	.fade-enter, .fade-leave-active {
		opacity: 0
	}
</style>
