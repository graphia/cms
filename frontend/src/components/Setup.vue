<template>

	<div class="row">
		<div class="col-md-9">
			<h1>Initial Setup</h1>

			<form @submit="create">

				<div class="form-group">
					<label for="username">Username</label>
					<input name="username" class="form-control" v-model="user.username"/>
				</div>

				<div class="form-group">
					<label for="name">Name</label>
					<input name="name" class="form-control" v-model="user.name"/>
				</div>

				<div class="form-group">
					<label for="email">Email</label>
					<input type="email" name="email" class="form-control" v-model="user.email"/>
				</div>

				<div class="form-group">
					<label for="password">Password</label>
					<input type="password" name="password" class="form-control" v-model="user.password"/>
				</div>

				<div class="form-group">
					<label for="confirm-password">Confirm Password</label>
					<input type="password" name="confirm-password" class="form-control" v-model="user.confirm_password"/>
				</div>

				<div class="form-group">
					<input type="submit" value="Create" class="btn btn-success">
				</div>

			</form>
		</div>
	</div>

</template>

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
				}
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