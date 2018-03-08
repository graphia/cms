<template>
	<div>
		<h1>New user</h1>

		<form @submit="create">
			<div class="form-group">
				<label for="username">Username</label>

				<input name="username"
					class="form-control"
					required
					minlength="3"
					maxlength="32"
					v-model="user.username"
					placeholder="milhouse.van.houten"
				/>
			</div>

			<div class="form-group">
				<label for="name">Name</label>

				<input name="name"
					class="form-control"
					v-model="user.name"
					placeholder="Milhouse van Houten"
					required
					minlength="3"
					maxlength="64"
				/>
			</div>

			<div class="form-group">
				<label for="email">Email</label>

				<input name="email"
					type="email"
					class="form-control"
					v-model="user.email"
					placeholder="milhouse.van.houten@k12.springfield.us"
					required
				/>
			</div>

			<div class="form-group">
				<label for="email">
					<input name="email" type="checkbox" v-model="user.admin"/>
					Administrator
				</label>

			</div>

			<div class="btn-toolbar">
				<input type="submit" value="Create user" class="btn btn-success"/>

				<router-link class="btn btn-secondary" :to="{name: 'user_settings', params: {id: user.id}}">
					Cancel
				</router-link>
			</div>
		</form>


	</div>
</template>


<script>

	import checkResponse from '../../../javascripts/response.js';
	import store from '../../../javascripts/store.js';
	import config from '../../../javascripts/config.js';

	class User {

		constructor(name, username, email, password, admin = false) {
			this.name = name;
			this.username = username;
			this.email = email;
			this.admin = admin;
		};

		async save() {
			const path = `${config.admin}/users`;

			let response = await fetch(path, {
				method: "POST",
				headers: store.state.auth.authHeader(),
				body: JSON.stringify({
					name: this.name,
					username: this.username,
					email: this.email,
					admin: this.admin
				})
			});
		}
	}

	export default {
		name: "NewUser",
		data() {
			return {
				user: new User
			}
		},
		methods: {
			async create(event) {
				event.preventDefault();

				let response = await this.user.save()



			}
		}
	};
</script>
