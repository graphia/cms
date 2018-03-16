<template>
	<div class="row">

		<div class="col-md-12">
			<h1>Welcome {{ this.user.name }}</h1>

			<p>
				Your account has been set up by an administrator but before
				you log in we need to configure your password.
			</p>

			<form id="activate-user" @submit="saveUser">

				<div class="form-group">
					<label class="form-control-label" for="user-password">Password</label>
					<input
						name="user-password"
						class="form-control"
						:class="{'is-valid': passwordsMatch && !confirmEmpty}"
						type="password"
						v-model="user.password"
						required="true"
						minlength="6"
						autocomplete="off"
					/>
				</div>

				<div class="form-group">
					<label class="form-control-label" for="user-confirm">Confirm Password</label>
					<input
						name="user-confirm"
						class="form-control"
						:class="{'is-valid': passwordsMatch && !confirmEmpty, 'is-invalid': !passwordsMatch && !confirmEmpty}"
						type="password"
						v-model="user.confirm"
						required="true"
						autocomplete="off"
					/>
				</div>

				<div class="form-group">
					<input
						type="submit"
						value="Update password"
						class="btn btn-primary"
						:disabled="anyEmpty || !passwordsMatch"
					/>
				</div>

			</form>
		</div>

	</div>
</template>

<script lang="babel">

	import checkResponse from '../../javascripts/response.js';

	class LimitedUser {

		constructor(params) {
			if (params) {
				this.username = params.username;
				this.name = params.name;
				this.confirmation_key = params.confirmation_key;
			};
			this.password = "";
			this.confirm = "";
		};

		passwordsMatch() {
			return this.password == this.confirm;
		};

		passwordEmpty() {
			return this.password == "";
		};

		confirmEmpty() {
			return this.confirm == "";
		};

		static async retrieve(ck) {
			const path = `/setup/activate/${ck}`;
			return await fetch(path);
		};

		async save(ck) {
			const path = `/setup/activate/${ck}`;
			return await fetch(path, {
				method: "PATCH",
				body: JSON.stringify({
					password: this.password
				})
			});
		};

	};

	export default {
		name: "ActivateUser",

		created() {
			this.getUser();
		},

		data() {
			return {
				user: new LimitedUser
			}
		},

		computed: {
			passwordsMatch() {
				return this.user.passwordsMatch();
			},
			confirmEmpty() {
				return this.user.confirmEmpty();
			},
			anyEmpty() {
				return (this.user.confirmEmpty() || this.user.passwordEmpty());
			}
		},

		methods: {
			async getUser() {
				const pk = this.$route.params.confirmation_key;

				let response = await LimitedUser.retrieve(pk);

				if (!checkResponse(response.status)) {
					console.error("failed to retrieve user with confirmation code", pk);
					return;
				};

				let json = await response.json();
				this.user = new LimitedUser(json);
			},
			async saveUser(event) {

				event.preventDefault();

				const pk = this.$route.params.confirmation_key;

				let response = await this.user.save(pk);

				if (!checkResponse(response.status)) {
					console.error("failed to activate user", pk);
					return;
				};

				this.$store.state.broadcast.addMessage(
					"success",
					"Your account has been activated",
					"Please log in using your new password",
					3
				);

				this.$store.state.auth.redirectToLogin();

			}
		}
	};
</script>