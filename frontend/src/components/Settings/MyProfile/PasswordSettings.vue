<template>
	<div class="password-update card border-warning mb-3">

		<div class="card-body">

			<h2 class="card-title">Update password</h2>

			<div class="alert alert-danger" v-if="anyErrors">
				{{ this.error }}
			</div>

			<form id="passwords" @submit="updatePassword">

				<div class="form-group">
					<label class="form-control-label" for="currentPassword">Current Password</label>
					<input
						name="currentPassword"
						class="form-control"
						type="password"
						v-model="currentPassword"
						required="true"
						autocomplete="off"
					/>
				</div>

				<div class="form-group">
					<label class="form-control-label" for="newPassword">New Password</label>
					<input
						name="newPassword"
						class="form-control"
						:class="{'is-valid': (confirmPasswordMatch && !confirmPasswordEmpty)}"
						type="password"
						v-model="newPassword"
						required="true"
						minlength="6"
						autocomplete="off"
					/>
				</div>

				<div class="form-group">
					<label class="form-control-label" for="confirmPassword">Confirm Password</label>
					<input
						name="confirmPassword"
						class="form-control"
						:class="{'is-valid': (confirmPasswordMatch && !confirmPasswordEmpty), 'is-invalid': (!confirmPasswordMatch && !confirmPasswordEmpty)}"
						type="password"
						v-model="confirmPassword"
						required="true"
						autocomplete="off"
					/>
				</div>

				<div class="form-group">
					<input
						type="submit"
						value="Update password"
						class="btn btn-warning"
						:disabled="anyEmpty || !confirmPasswordMatch"
					/>
				</div>
			</form>


		</div>
	</div>
</template>

<script>
	import Accessors from '../../Mixins/accessors';
	import checkResponse from "../../../javascripts/response.js";

	export default {
		name: "PasswordSettings",
		data() {
			return {
				newPassword: "",
				currentPassword: "",
				confirmPassword: "",
				error: ""
			};
		},
		computed: {
			confirmPasswordMatch() {
				return (this.newPassword == this.confirmPassword);
			},
			confirmPasswordEmpty() {
				return (this.confirmPassword == "");
			},
			anyErrors() {
				return this.error != "";
			},
			anyEmpty() {
				return [this.newPassword, this.currentPassword, this.confirmPassword]
					.some(input => input === "");
			}
		},
		methods: {
			async updatePassword(event) {
				this.error = "";

				event.preventDefault();
				let response = await this.currentUser.updatePassword(
					this.currentPassword,
					this.newPassword
				);

				if (response.status == 400) {
					// bad request, show the error
					let json = await response.json();
					this.error = json.message;
					return

				} else if (!checkResponse(response.status)) {
					// any other errors
					throw {reason: "Couldn't get user info", code: response.status}
				};

				this.$store.state.broadcast.addMessage(
					"success",
					"Password updated",
					"Next time you log in, use your new password",
					3
				);

				this.reset();
			},
			reset() {
				this.newPassword = "";
				this.currentPassword = "";
				this.confirmPassword = "";
			}
		},
		mixins: [Accessors]

	}
</script>

