<template>
	<div class="col col-md-6" v-title="title">

		<Breadcrumbs :levels="breadcrumbs"/>

		<h1>{{ title }}</h1>

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
				<label for="email">Email address</label>

				<input name="email"
					type="email"
					class="form-control"
					v-model="user.email"
					placeholder="milhouse.van.houten@k12.springfield.us"
					required
				/>
			</div>

			<div class="form-group">
				<label for="admin">
					<input name="admin" type="checkbox" v-model="user.admin"/>
					Administrator
				</label>

			</div>

			<div class="btn-toolbar" role="group">
				<input type="submit" value="Create user" class="btn btn-success"/>

				<router-link class="btn btn-secondary ml-2" :to="{name: 'user_settings', params: {id: user.id}}">
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

	import Breadcrumbs from '../../Utilities/Breadcrumbs';
	import CMSBreadcrumb from '../../../javascripts/models/breadcrumb.js';

	class User {

		constructor(name, username, email, admin = false) {
			this.name = name;
			this.username = username;
			this.email = email;
			this.admin = admin;
		};

		async save() {
			const path = `${config.admin}/users`;

			return fetch(path, {
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
				user: new User,
				title: "New user"
			}
		},
		methods: {
			async create(event) {
				event.preventDefault();

				let response = await this.user.save();

				if (!checkResponse(response.status)) {
					console.error("Failed to create user", response);
					return;
				};

				this.$store.state.broadcast.addMessage(
					"success",
					"User Created",
					`${this.user.name} will receive an email with instructions on how to log in`,
					3
				);

				this.$router.push({name: 'user_settings'});

			}
		},
		computed: {
			breadcrumbs() {
				return [
					new CMSBreadcrumb("Users", "user_settings"),
					new CMSBreadcrumb("New", "user_new")
				];
			},
		},
		components: {
			Breadcrumbs
		}
	};
</script>
