<template>
	<div class="document-title form-group" v-bind:class="{'has-danger': !valid}">

		<label for="title">Title</label>

		<input
			id="document-title"
			name="title"
			class="form-control"
			v-model="document.title"
			minlength="2"
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
		name: "TitleField",
		data() {
			return {
				element: null,
				valid: true,
				validationMessage: null
			};
		},
		computed: {
			document() {
				return this.$store.state.activeDocument;
			}
		},
		methods: {
			validate() {
				// make the parent validate the whole form to control
				// display of the submit button
				this.$bus.$emit("checkMetadata");

				if (!this.element) {
					this.element = document.getElementById("document-title");
				};

				this.valid = this.element.checkValidity();
				this.validationMessage = this.element.validationMessage;
			}
		}
	};
</script>