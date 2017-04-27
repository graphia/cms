<template>
	<div>
		<form @submit="login">
			<label for="username">Username</label>
			<input name="username" v-model="username"/>
			<label for="password">Password</label>
			<input type="password" name="password" v-model="password"/>

			<input type="submit" value="Log in"/>
		</form>
	</div>
</template>

<script lang="babel">
	export default {
		name: "Login",
		data() {
			return {
				username: "",
				password: ""
			};
		},
		methods: {
			async login(event) {
				event.preventDefault();
				console.log("clicked!");

				let response = await fetch(`${this.$config.auth}/login`, {
					method: "POST",
					mode: "cors",
					body: JSON.stringify({
						username: this.username,
						password: this.password
					})
				});

				if (response.status !== 200) {
					console.error('Oops, there was a problem', response.status);
					return
				}

				let json = await response.json();
				localStorage.setItem('token', json.token);

				// token saved, redirect to somewhere useful!


			}
		}
	};
</script>