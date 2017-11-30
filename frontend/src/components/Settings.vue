<template>
	<div>
		<h1>Settings</h1>

		<div class="personal-details">
			<h4>Personal Details</h4>


			<form @submit="updateUser">


				<div class="form-group">
					<label class="form-control-label" for="name">Name</label>
					<input
						name="name"
						class="form-control"
						placeholder="Charles Montgomery Burns"
						v-model="user.name"
						required="true"
						readonly="true"
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
						class="btn btn-primary"
					/>
				</div>
			</form>

		</div>

	</div>
</template>

<script>
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
		}
	};
</script>
