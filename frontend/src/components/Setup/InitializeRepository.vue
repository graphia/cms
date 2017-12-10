<template>
	<div class="row">

		<div class="col col-md-9">
			<div class="alert alert-warning">
				<p>
					<octicon :icon-name="alert"></octicon>

					<strong>No repository found.</strong> If this is a fresh install you can complete the process
					by initialising the brand new repository.
				</p>
			</div>

			<button @click="initializeRepo" class="btn btn-primary btn-lg">Initialise Repository</button>
		</div>
	</div>
</template>

<script lang="babel">

	import config from '../../javascripts/config.js';
	import checkResponse from '../../javascripts/response.js';

	export default {
		name: "SetupInitialiseRepository",
		methods: {
			async initializeRepo(event) {
				event.preventDefault();


				let path = `${config.api}/setup/initialize_repository`;

				try {
					let response = await fetch(path, {
						method: "POST",
						mode: "cors",
						headers: this.$store.state.auth.authHeader()
					});

					if (!checkResponse(response.status)) {
						return;
					}

					let json = await response.json();

					this.$store.state.broadcast.addMessage(
						"success",
						"Repository initialised", "You have created a new place to store your documents, start by adding a new directory",
						10
					);

					this.redirectToHome();

			}
				catch (error) {
					console.error(error);
				}

			},
			redirectToHome() {
				this.$router.push({name: 'home'});
			}
		}
	};
</script>
