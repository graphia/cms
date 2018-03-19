<template>
	<div class="col col-md-6" v-title="title">

		<Breadcrumbs :levels="breadcrumbs"/>

		<h1>{{ title }}</h1>

		<form @submit="create">

			<div class="form-group" v-if="hasErrors('Base')">
				<div class="alert alert-danger">
					This record cannot be saved because either the username or email address
					are already in use
				</div>
			</div>

			<div class="form-group">
				<label for="username">Username</label>

				<input name="username"
					class="form-control"
					required
					minlength="3"
					maxlength="32"
					v-model="user.username"
					placeholder="milhouse.van.houten"
					:class="{'is-invalid': hasErrors('Username')}"
				/>

				<div class="form-control-feedback invalid-feedback" v-if="hasErrors('Username')">
					{{ errorMessage('Username') }}
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
				title: "New user",
				errors: {}
			}
		},
		methods: {
			async create(event) {
				event.preventDefault();

				let response = await this.user.save();

				let json = await response.json();

				if (!checkResponse(response.status)) {
					this.errors = json;
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
					new CMSBreadcrumb("New", "user_new")
				];
			},
		},
		components: {
			Breadcrumbs
		}
	};
</script>
