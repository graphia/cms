<template>
	<ul class="col-md-3 nav nav-pills flex-column p-3">
		<li class="nav-item">
			<router-link class="nav-link" :to="{name: 'my_profile'}" linkActiveClass="active">
				My Profile
			</router-link>
		</li>
		<li class="nav-item">
			<router-link class="nav-link" :to="{name: 'ssh_key_settings'}" linkActiveClass="active">
				SSH Keys
			</router-link>
		</li>
		<li class="nav-item">
			<router-link class="nav-link" :to="{name: 'user_settings'}" linkActiveClass="active">
				Users
			</router-link>
		</li>
	</ul>
</template>

<script>
	export default {
		name: "SettingsNavigation",
	}
</script>
