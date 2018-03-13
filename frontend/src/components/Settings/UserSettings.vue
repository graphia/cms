<template>
	<div v-title="title">

		<Breadcrumbs :levels="breadcrumbs"/>

		<div class="row">
			<SettingsNavigation/>

			<div class="col-md-9">

				<div class="row">
					<div class="col text-right p-2 m-2">
						<router-link class="btn btn-secondary" :to="{name: 'user_new'}">
							Create new user
						</router-link>
					</div>
				</div>

				<div class="row user-list">

					<div class="col px-2" v-for="(user, i) in users" :key="i">
						<div class="card user my-2" :id="`user-${user.id}`">

							<div class="card-header">
								{{ user.name }}

								<span class="badge badge-info" v-if="user.admin">
									Admin
								</span>

								<span class="badge badge-warning" v-if="!user.active">
									Deactivated
								</span>
							</div>

							<div class="card-body">
								<dl>
									<dt>Username:</dt>
									<dd>{{ user.username }}</dd>

									<dt>Email Address:</dt>
									<dd>{{ user.email }}</dd>
								</dl>
							</div>

							<div class="card-footer">
								<router-link class="btn btn-secondary" :to="{name: 'user_edit', params: {username: user.username}}">
									Edit
								</router-link>
							</div>
						</div>
					</div>

				</div>

			</div>
		</div>
	</div>

</template>


<script lang="babel">

	import SettingsNavigation from "./Navigation";
	import Breadcrumbs from '../Utilities/Breadcrumbs';
	import PasswordSettings from './MyProfile/PasswordSettings';


	import checkResponse from '../../javascripts/response.js';
	import store from '../../javascripts/store.js';
	import config from '../../javascripts/config.js';
	import CMSUser from '../../javascripts/models/user.js';
	import CMSBreadcrumb from '../../javascripts/models/breadcrumb.js';

	class UserList {
		static async all() {
			const path = `${config.api}/users`;
			let response = await fetch(path, {headers: store.state.auth.authHeader()})

			if (!checkResponse(response.status)) {
				console.error("User list cannot be retrieved", response);
				return;
			};

			return await response.json();

		};
	};

	export default {
		name: "UserSettings",
		async created() {
			this.users = await UserList.all();
		},
		data() {
			return {
				title: "Users",
				users: []
			};
		},
		computed: {
			breadcrumbs() {
				return [new CMSBreadcrumb("Users", "user_settings")];
			},
		},
		components: {
			SettingsNavigation,
			Breadcrumbs,
			PasswordSettings
		}
	};
</script>
