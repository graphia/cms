<template>
	<div id="application">

		<header>
			<!-- Primary Navigation Start -->
			<nav class="navbar navbar-expand-md navbar-dark bg-dark">

				<router-link class="navbar-brand" :to="{name: 'home'}">{{ this.siteTitle }}</router-link>


				<button v-if="user" class="navbar-toggler navbar-toggler-right hidden-md-up" type="button" data-toggle="collapse" data-target="#primary" aria-label="Toggle navigation">
					<span class="navbar-toggler-icon"></span>
				</button>

				<div v-if="user" id="primary" class="collapse navbar-collapse">
					<ul class="navbar-nav mr-auto">

						<li>
							<router-link :to="{name: 'home'}" exact class="nav-link home-link">Home</router-link>
						</li>

						<li v-for="(directory, i) in directories" :key="i">
							<router-link :to="{name: 'directory_index', params: {directory: directory.path}}" class="nav-link directory-link">
								{{ directory.title || directory.path | capitalize }}
							</router-link>
						</li>

					</ul>

					<ul class="navbar-nav">
						<li id="user-dropdown" class="nav-item dropdown">
							<a class="nav-link dropdown-toggle" href="#" id="user-menu" data-toggle="dropdown" aria-haspopup="true" aria-expanded="false">
								{{ user.persistedName }}
							</a>
							<div class="dropdown-menu" aria-labelledby="user-menu">
								<router-link :to="{name: 'my_profile'}" class="dropdown-item">
									Settings
								</router-link>
								<a class="dropdown-item logout" href="logout" @click="logout">Logout</a>
							</div>
						</li>
					</ul>
				</div>

			</nav>
			<!-- Primary Navigation End -->
		</header>

		<!-- Router View Container Start -->
		<div class="main container-fluid">
			<Broadcast/>
			<transition name="fade">
				<router-view/>
			</transition>
		</div>
		<!-- Router View Container End -->

		<footer>
			<div class="copyright m-2">
				<a href="https://www.graphia.co.uk">Graphia CMS&trade;</a>
			</div>
		</footer>

	</div>
</template>

<script lang="babel">
	import config from '../javascripts/config.js';

	import checkResponse from '../javascripts/response.js';
	import CMSAuth from '../javascripts/auth.js';
	import Broadcast from '../components/Broadcast';
	import Accessors from './Mixins/accessors';

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

				console.info("token is present and has not expired, renewing");
				this.$store.state.auth.renew();

				["getTopLevelDirectories"]
					.map(func => {
						this.$store.dispatch(func);
					});

				["refreshTranslationInfo", "refreshRepositoryInfo", "refreshServerInfo"]
					.map(func => {
						this.$store.commit(func);
					});

				// load user data if it's not present from a fresh login
				if (!this.$store.state.user) {
					this.$store.commit("loadUser");
				};

			}
			catch(err) {
				// Token rejected for renewal

				let allow = CMSAuth.unblockedPageCheck(this.$route.path);

				if (!allow) {
					console.warn("Token not valid", err);
					this.$store.state.auth.redirectToLogin();
				};
			};

		},

		computed: {
			user() {
				return this.$store.state.user;
			},
			siteTitle() {
				return (this.$store.state.server.serverInfo.title || "Graphia CMS");
			}
		},
		methods: {
			logout(event) {
				event.preventDefault();
				this.$store.state.auth.logout();
			}
		},
		components: {
			Broadcast
		},
		mixins: [Accessors],

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
		opacity: 0;
	}
</style>
