<template>

	<div class="row">

		<div class="col-sm-8 offset-sm-2">

			<div class="card mt-4">

				<h3 class="card-header">Create an administrator</h3>

				<!-- Setup form start -->
				<div class="card-body">

					<form @submit="create">

						<div class="form-group">
							<label for="username">Username</label>
							<input name="username" class="form-control" v-model="user.username" minlength="3" maxlength="32" required/>
						</div>

						<div class="form-group">
							<label for="name">Full Name</label>
							<input name="name" class="form-control" v-model="user.name" required minlength="3" maxlength="64"/>
						</div>

						<div class="form-group">
							<label for="email">Email Address</label>
							<input type="email" name="email" class="form-control" v-model="user.email" required/>
						</div>

						<div class="form-group">
							<label for="password">Password</label>
							<input type="password" name="password" class="form-control" v-model="user.password" required minlength="6"/>
						</div>

						<div class="form-group confirm-password-group">
							<label for="confirm-password">Confirm Password</label>
							<input
								type="password"
								name="confirm-password"
								class="form-control"
								:class="[{'is-invalid': !passwordsMatch}]"
								v-model="user.confirm_password"
								required
							/>

							<div v-if="!passwordsMatch" class="password-match-feedback form-control-feedback-message">Password and confirmation must match</div>
						</div>

						<div class="form-group" >
							<input type="submit" value="Create" class="btn btn-success">
						</div>

					</form>

				</div>
				<!-- Setup form end -->

			</div>

		</div>

	</div>
</template>

<script>
	import CMSAuth from '../../javascripts/auth.js';

	export default {
		name: "Setup",
		data() {
			return {
				user: {
					username: null,
					name: null,
					password: null,
					confirm_password: null,
					email: null
				}
			}
		},
		computed: {
			passwordsMatch() {
				return (this.user.password == this.user.confirm_password);
			}
		},
		methods: {
			async create(event) {
				event.preventDefault();
				await CMSAuth.createInitialUser(this.user);
				this.$store.state.broadcast.addMessage("success", "Welcome", "Administrator created", 3);
				this.$router.push({name: 'login'});
			}
		}
	}
</script>
