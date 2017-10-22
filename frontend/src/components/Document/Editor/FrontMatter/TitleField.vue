<template>
	<div class="document-title form-group">

		<label for="title" class="form-control-label">Title</label>

		<input
			id="document-title"
			name="title"
			class="form-control"
			v-model="document.title"
			minlength="2"
			required="true"
			v-on:keyup="validate"
			:class="{'is-invalid': !valid}"
		/>

		<div class="form-control-feedback invalid-feedback" v-if="validationMessage">
			{{ this.validationMessage }}
		</div>
	</div>
</template>

<script lang="babel">

	import Accessors from '../../../Mixins/accessors';

	export default {
		name: "TitleField",
		data() {
			return {
				element: null,
				valid: true,
				validationMessage: null
			};
		},
		mixins: [Accessors],
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