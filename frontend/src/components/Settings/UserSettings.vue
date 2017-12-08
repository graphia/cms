<template>
	<div>

		<Breadcrumbs :levels="breadcrumbs"/>

		<div class="row">
			<SettingsNavigation/>

			<div class="col-md-9 personal-details">
				<h1>Personal Details</h1>

				<form id="personal-details" @submit="updateUser">

					<div class="form-group">
						<label class="form-control-label" for="username">Username</label>
						<input
							name="username"
							class="form-control"
							type="username"
							placeholder="monty.burns@springfieldpower.com"
							v-model="user.username"
							required="true"
							disabled="true"
						/>
					</div>

					<div class="form-group">
						<label class="form-control-label" for="email">Email</label>
						<input
							name="email"
							class="form-control"
							type="email"
							placeholder="monty.burns@springfieldpower.com"
							v-model="user.email"
							required="true"
							disabled="true"
						/>
					</div>

					<div class="form-group">
						<label class="form-control-label" for="name">Name</label>
						<input
							name="name"
							class="form-control"
							placeholder="Charles Montgomery Burns"
							v-model="user.name"
							required="true"
							minlength="3"
							maxlength="64"
						/>
					</div>

					<div class="form-group">
						<input
							type="submit"
							value="Update my details"
							class="btn btn-success"
							:disabled="!user.updated()"

						/>
					</div>
				</form>


				<PasswordSettings/>


			</div>

		</div>

	</div>
</template>

<script>

	import SettingsNavigation from "./Navigation";
	import Breadcrumbs from '../Utilities/Breadcrumbs';
	import PasswordSettings from './UserSettings/PasswordSettings';

	import CMSUser from '../../javascripts/models/user.js';
	import CMSBreadcrumb from '../../javascripts/models/breadcrumb.js';

	export default {
		name: "Settings",
		data() {
			return {
				 // so we can render the form bound to something
				emptyUser: new CMSUser
			};
		},
		computed: {
			user() {
				return (this.$store.state.user || this.emptyUser);
			},

			breadcrumbs() {
				return [new CMSBreadcrumb("User settings", "user_settings")];
			}

		},
		methods: {
			async updateUser(event) {
				event.preventDefault();
				console.debug("form submitted!");
				await this.user.save();
				this.user.refresh();
			}
		},
		components: {
			SettingsNavigation,
			Breadcrumbs,
			PasswordSettings
		}

	};
</script>
