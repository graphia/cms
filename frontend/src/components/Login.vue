<template>
	<div class="row" v-title="title">

		<div class="col-sm-6 offset-sm-3">

			<div class="card mt-4">

				<h3 class="card-header">Log in</h3>

				<!-- Login form start -->
				<div class="card-body">

					<form @submit="login" class="login-form">

						<div class="form-group">
							<label for="username">Username</label>
							<input class="form-control" type="text" name="username" v-model="username" required/>
						</div>
						<div class="form-group">
							<label for="password">Password</label>
							<input class="form-control" type="password" name="password" v-model="password" required/>
						</div>

						<div class="form-group">
							<input type="submit" value="Log in" class="btn btn-primary"/>
						</div>
					</form>

				</div>
				<!-- Login form end -->

			</div>

		</div>
	</div>
</template>

<style lang="scss">
	.card-body .form-group:last-child {
		margin-bottom: 0em;
	}
</style>

<script lang="babel">
	import CMSAuth from '../javascripts/auth.js'
	export default {
		name: "Login",
		data() {
			return {
				username: "",
				password: "",
				title: "Graphia CMS: Login"
			};
		},
		created() {
			this.redirectToSetup();
		},
		methods: {
			async login(event) {
				event.preventDefault();

				let success = await this.$store.state.auth.login(this.username, this.password);

				if (!success) {
					this.$store.state.broadcast.addMessage("danger", "Oops", "Invalid credentials", 5);
					return;
				}

				this.$store.state.broadcast.addMessage("success", "Welcome", "You have logged in successfully", 3);

				// if we've stored the original destination (globally), use it and clear it
				if (window.originalDestination) {
					this.$router.push(window.originalDestination);
					delete window.originalDestination;
					return;
				};

				this.$router.push({name: 'home'});

			},
			async redirectToSetup() {

				// check if there are any users
				let doSetup = await CMSAuth.doInitialSetup();

				// if there are, abort!
				if (!doSetup) {
					return;
				}

				// if there aren't, start the setup wizard
				this.$router.push({
					name: 'initial_setup'
				});
			}
		}
	};
</script>
