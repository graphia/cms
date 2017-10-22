<template>

	<div class="commit-message form-group">

		<label class="form-control-label" for="commit-message">
			Commit Message
		</label>

		<textarea
			id="document-commit-message"
			name="commit-message"
			class="form-control"
			v-model="commit.message"
			minlength="5"
			required="true"
			v-on:keyup="validate"
			v-bind:class="{'is-invalid': !valid}"
		/>

		<div class="form-control-feedback invalid-feedback" v-if="validationMessage">
			{{ this.validationMessage }}
		</div>

	</div>

</template>

<script lang="babel">
	import Accessors from '../../Mixins/accessors';

	export default {
		name: "CommitMessageField",
		data() {
			return {
				element: null,
				valid: true,
				validationMessage: null
			};
		},
		methods: {
			validate() {
				// make the parent validate the whole form to control
				// display of the submit button
				this.$bus.$emit("checkMetadata");

				if (!this.element) {
					this.element = document.getElementById("document-commit-message");
				};

				this.valid = this.element.checkValidity();
				this.validationMessage = this.element.validationMessage;
			}
		},
		mixins: [Accessors]

	}
</script>