<template>
	<div class="row">

		<SettingsNavigation/>

		<div class="col-md-9 personal-details">
			<h1>Personal Details</h1>

			<form @submit="updateUser">

				<div class="form-group">
					<label class="form-control-label" for="username">Username</label>
					<input
						name="username"
						class="form-control"
						type="username"
						placeholder="monty.burns@springfieldpower.com"
						v-model="user.username"
						required="true"
						disabled="true"
					/>
				</div>

				<div class="form-group">
					<label class="form-control-label" for="email">Email</label>
					<input
						name="email"
						class="form-control"
						type="email"
						placeholder="monty.burns@springfieldpower.com"
						v-model="user.email"
						required="true"
						disabled="true"
					/>
				</div>

				<div class="form-group">
					<label class="form-control-label" for="name">Name</label>
					<input
						name="name"
						class="form-control"
						placeholder="Charles Montgomery Burns"
						v-model="user.name"
						required="true"
					/>
				</div>

				<div class="form-group">
					<input
						type="submit"
						value="Update my details"
						class="btn btn-success"
						:disabled="!user.updated()"

					/>
				</div>
			</form>

		</div>

	</div>
</template>

<script>

	import SettingsNavigation from "./Navigation";

	export default {
		name: "Settings",
		computed: {
			user() {
				return this.$store.state.user;
			}
		},
		methods: {
			updateUser(event) {
				event.preventDefault();
				console.debug("form submitted!");
			}
		},
		components: {
			SettingsNavigation
		}

	};
</script>
