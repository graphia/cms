<template>
	<div class="directory-title form-group">

		<label class="form-control-label" for="title">Title</label>
		<input
			:id="elementID"
			name="title"
			type="text"
			class="form-control"
			placeholder="Operating Procedures"
			v-model="activeDirectory.title"
			required="true"
			minlength=2
			autocomplete="off"
			v-on:keyup="validate"
			:class="{'is-invalid': !valid}"
		/>

		<div class="form-control-feedback invalid-feedback" v-if="validationMessage">
			{{ this.validationMessage }}
		</div>
	</div>
</template>

<script lang="babel">

	import Accessors from '../../Mixins/accessors';

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
		mixins: [Accessors],
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