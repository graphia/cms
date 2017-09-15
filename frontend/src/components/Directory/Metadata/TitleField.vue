<template>
	<div class="directory-title form-group" :class="{'has-danger': !valid}">

		<label for="title">Title</label>
		<input
			:id="elementID"
			name="title"
			type="text"
			class="form-control"
			placeholder="Operating Procedures"
			v-model="directory.title"
			required="true"
			minlength=2
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
				validationMessage: null,
				elementID: "directory-title"
			};
		},
		computed: {
			directory() {
				return this.$store.state.activeDirectory;
			}
		},
		methods: {
			validate() {
				this.$bus.$emit("checkMetadata");

				if (!this.element) {
					this.element = document.getElementById(this.elementID);
				};

				this.valid = this.element.checkValidity();
				this.validationMessage = this.element.validationMessage;
			}
		}
	};
</script>