<template>
	<div v-title="title">

		<Breadcrumbs :levels="breadcrumbs"/>

		<div class="row">
			<SettingsNavigation/>

			<div class="col-md-9 user-list">
				<div v-for="(user, i) in users" :key="i">

					<h1>{{ user.name }}</h1>
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
				return [new CMSBreadcrumb("User List", "user_settings")];
			},
		},
		components: {
			SettingsNavigation,
			Breadcrumbs,
			PasswordSettings
		}
	};
</script>
