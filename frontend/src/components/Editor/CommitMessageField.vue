<template>

	<div class="commit-message form-group" v-bind:class="{'has-danger': !valid}">

		<label
			class="form-control-label"
			for="commit-message"
		>
			Commit Message
		</label>

		<textarea
			id="new-document-commit-message"
			name="commit-message"
			class="form-control"
			v-model="commit.message"
			minlength="5"
			required="true"
			v-on:keyup="validate"
		/>

		<div class="form-control-feedback" v-if="validationMessage">
			{{ this.validationMessage }}
		</div>

	</div>

</template>

<script lang="babel">
	export default {
		name: "CommitMessage",
		data() {
			return {
				element: null,
				valid: true,
				validationMessage: null
			};
		},
		computed: {
			commit() {
				return this.$store.state.commit;
			}
		},

		methods: {
			validate() {
				// make the parent validate the whole form to control
				// display of the submit button
				this.$bus.$emit("checkMetadata");

				if (!this.element) {
					this.element = document.getElementById("new-document-commit-message");
				};

				this.valid = this.element.checkValidity();
				this.validationMessage = this.element.validationMessage;
			}
		}

	}
</script>