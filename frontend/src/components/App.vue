<template>
	<div id="application">

		<!-- Primary Navigation Start -->
		<nav class="navbar navbar-toggleable-md navbar-inverse static-top bg-inverse">
			<button class="navbar-toggler navbar-toggler-right hidden-lg-up" type="button" data-toggle="collapse" data-target="#primary" aria-label="Toggle navigation">
				<span class="navbar-toggler-icon"></span>
			</button>

			<a class="navbar-brand" href="#">Graphia CMS</a>

			<div id="primary" class="collapse navbar-collapse">
				<ul class="navbar-nav mr-auto">
					<li class="nav-item active">Files</li>

					<router-link :to="{name: 'document_index', params: {directory: 'documents'}}" class="nav-link">
						Documents
					</router-link>

					<li><a class="nav-link" href="#">History</a></li>
					<li><a class="nav-link" href="#">Admin</a></li>

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
	import CMSAuth from '../javascripts/auth.js';
	import Broadcast from '../components/Broadcast';

	export default {
		name: "GraphiaCMS",
		created() {
			try {

				if (!this.$store.state.auth.tokenExpiry) {
					throw 'Token missing';
				}

				if (this.$store.state.auth.expiry < Date.now()) {
					throw 'Token expired';
				}

				console.debug("token is present and has not expired, renewing");
				this.$store.state.auth.renew();

			}
			catch(err) {
				// Token rejected for renewal
				console.warn(err);
				this.$store.state.auth.redirectToLogin();
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