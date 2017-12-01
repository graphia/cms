<template>
	<div class="row">

		<SettingsNavigation/>

		<div class="col-md-9 personal-details">
			<h1>Personal Details</h1>

			<form @submit="updateUser">

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
					<label class="form-control-label" for="email">Email</label>
					<input
						email="name"
						class="form-control"
						type="email"
						placeholder="monty.burns@springfieldpower.com"
						v-model="user.email"
						required="true"
						readonly="true"
					/>
				</div>

				<div class="form-group">
					<input
						type="submit"
						value="Update my details"
						class="btn btn-success"

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
		created() {
			this.orignalUser = this.user;
		},
		data() {
			return {
				orignalUser: null,
			};
		},
		computed: {
			user() {
				return this.$store.state.user;
			},
			userModified() {
				return (this.user === this.originalUser);
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
