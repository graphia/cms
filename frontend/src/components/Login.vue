<template>
	<div class="row">

		<div class="col-sm-6 offset-sm-3">

			<div class="card mt-4">

				<div class="card-header">
					<ul class="nav nav-tabs card-header-tabs">
						<li class="nav-item">
							<a class="nav-link active" href="#">Login</a>
						</li>
						<li class="nav-item">
							<a class="nav-link" href="#">Sign up</a>
						</li>
					</ul>
				</div>

				<!-- Login form start -->
				<div class="card-block">

					<form @submit="login">

						<div class="form-group">
							<label for="username">Username</label>
							<input class="form-control" name="username" v-model="username" required/>
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

<script lang="babel">
	export default {
		name: "Login",
		data() {
			return {
				username: "",
				password: ""
			};
		},
		methods: {
			async login(event) {
				event.preventDefault();
				console.log("clicked!");

				await this.$store.state.auth.login(this.username, this.password);

				// TODO if they'd attempted to navigate to a page
				// we should store it and send them there.
				this.$router.push({name: 'home'});

			}
		}
	};
</script>