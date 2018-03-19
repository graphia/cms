<template>
	<div class="col col-md-6" v-title="title">

		<div v-if="initializing">
			Loading...
		</div>

		<div v-else-if="user.found">

			<Breadcrumbs :levels="breadcrumbs"/>

			<h1>{{ title }}</h1>

			<form @submit="update">

				<div class="form-group" v-if="hasErrors('Base')">
					<div class="alert alert-danger">
						This record cannot be saved because either the username or email address
						are already in use
					</div>
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
						:class="{'is-invalid': hasErrors('Name')}"
					/>

					<div class="form-control-feedback invalid-feedback" v-if="hasErrors('Name')">
						{{ errorMessage('Name') }}
					</div>
				</div>

				<div class="form-group">
					<label for="email">Email address</label>

					<input name="email"
						type="email"
						class="form-control"
						v-model="user.email"
						placeholder="milhouse.van.houten@k12.springfield.us"
						:class="{'is-invalid': hasErrors('Email')}"
						required
					/>

					<div class="form-control-feedback invalid-feedback" v-if="hasErrors('Email')">
						{{ errorMessage('Email') }}
					</div>
				</div>

				<div class="form-group">
					<label for="admin">
						<input name="admin" type="checkbox" v-model="user.admin"/>
						Administrator
					</label>

				</div>

				<div class="form-group">
					<label for="active">
						<input name="active" type="checkbox" v-model="user.active"/>
						Active
					</label>

				</div>

				<div class="btn-toolbar" role="group">
					<input type="submit" value="Update user" class="btn btn-success"/>

					<router-link class="btn btn-secondary mx-2" :to="{name: 'user_settings', params: {id: user.id}}">
						Cancel
					</router-link>

					<button class="btn btn-danger" @click="deleteUser">
						Delete
					</button>
				</div>
			</form>

		</div>

		<div v-else>
			<h1>404</h1>
			User {{ this.$route.params.username }} cannot be found.
		</div>

	</div>
</template>


<script>

	import checkResponse from '../../../javascripts/response.js';
	import store from '../../../javascripts/store.js';
	import config from '../../../javascripts/config.js';
	import Breadcrumbs from '../../Utilities/Breadcrumbs';
	import CMSBreadcrumb from '../../../javascripts/models/breadcrumb.js';

	class User {

		constructor(name, username, email, admin=false, active=true, found=false) {
			this.name = name;
			this.username = username;
			this.email = email;
			this.admin = admin;
			this.active = active;
			this.found = found;
		};

		static async find(username) {

			const path = `${config.api}/users/${username}`;
			let response = await fetch(path, {headers: store.state.auth.authHeader()});

			if (!checkResponse(response.status)) {
				console.error(`User '${username}' cannot be retrieved`);

				if (response.status == 404) {
					console.error("user not found");
					return response.status;
				};

				return;
			}
			let json = await response.json();

			return new User(json.name, json.username, json.email, json.admin, json.active, true);
		};

		async save() {
			const path = `${config.admin}/users/${this.username}`;

			return fetch(path, {
				method: "PATCH",
				headers: store.state.auth.authHeader(),
				body: JSON.stringify({
					name: this.name,
					username: this.username,
					email: this.email,
					admin: this.admin,
					active: this.active
				})
			});

		};

		async delete() {
			const path = `${config.admin}/users/${this.username}`;
			return fetch(path, {
				method: "DELETE",
				headers: store.state.auth.authHeader()
			});
		};
	};

	export default {
		name: "EditUser",
		async created() {
			this.findUser();
		},
		data() {
			return {
				user: new User,
				title: "Edit user",
				errors: {},
				initializing: true,
				userFound: false
			};
		},
		methods: {
			async findUser() {
				const username = this.$route.params.username;
				this.user = await User.find(username);

				this.initializing = false;
			},
			async update(event) {
				event.preventDefault();

				let response = await this.user.save();
				let json = await response.json();

				if (!checkResponse(response.status)) {
					this.errors = json;
					console.error("Failed to modify user", response);
					return;
				};

				this.$store.state.broadcast.addMessage(
					"success",
					`User ${this.user.name} updated`,
					"",
					3
				);
				this.$router.push({name: 'user_settings'});
			},
			async deleteUser(event) {
				event.preventDefault();

				let response = await this.user.delete();
				let json = await response.json();

				if (!checkResponse(response.status)) {
					console.error("Failed to delete user", response);

					this.$store.state.broadcast.addMessage(
						"danger",
						`User ${this.user.name} not deleted`,
						json.message,
						3
					);

					return;
				};

				this.$store.state.broadcast.addMessage(
					"success",
					`User ${this.user.name} deleted`,
					"",
					3
				);

				this.$router.push({name: 'user_settings'});

			},

			hasErrors(field) {
				return !!this.errors[field];
			},
			errorMessage(field) {
				return this.errors[field];
			}
		},
		computed: {
			breadcrumbs() {
				return [
					new CMSBreadcrumb("Users", "user_settings"),
					new CMSBreadcrumb("Edit", "user_edit", {username: this.user.name})
				];
			},
		},
		components: {
			Breadcrumbs
		}
	};
</script>
