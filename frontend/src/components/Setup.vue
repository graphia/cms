<template>

	<div class="row">
		<div class="col-md-9">
			<h1>Initial Setup</h1>

			<form @submit="create">

				<div class="form-group">
					<label for="username">Username</label>
					<input name="username" class="form-control" v-model="user.username" minlength="4" required/>
				</div>

				<div class="form-group">
					<label for="name">Name</label>
					<input name="name" class="form-control" v-model="user.name" required/>
				</div>

				<div class="form-group">
					<label for="email">Email</label>
					<input type="email" name="email" class="form-control" v-model="user.email" required/>
				</div>

				<div class="form-group">
					<label for="password">Password</label>
					<input type="password" name="password" class="form-control" v-model="user.password" required/>
				</div>

				<div class="form-group confirm-password-group" v-bind:class="[{'has-danger passwords-do-not-match': !passwordsMatch}]">
					<label for="confirm-password">Confirm Password</label>
					<input type="password" name="confirm-password" class="form-control" v-model="user.confirm_password" required/>

					<div class="form-control-feedback passwords-do-not-match-message">Password and confirmation must match</div>
				</div>

				<div class="form-group" >
					<input type="submit" value="Create initial user" class="btn btn-success">
				</div>

			</form>
		</div>
	</div>

</template>

<style lang="scss">
	.passwords-do-not-match-message {
		display: none;
	}
	.passwords-do-not-match > .passwords-do-not-match-message {
		display: block;
	}
</style>

<script>
	import CMSAuth from '../javascripts/auth.js';

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
				},
				errorClass: 'has-danger'
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
				// TODO flash message
				this.$router.push({name: 'login'});
			}
		}
	}
</script>