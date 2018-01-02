<template>
	<div class="document-date form-group">
		<label for="date" class="form-control-label">Date</label>

		<input
			id="document-date"
			type="date"
			name="date"
			class="form-control"
			v-model="document.date"
			v-on:keyup="validate"
			required="true"
		/>
	</div>

</template>

<script lang="babel">
	import Accessors from '../../../Mixins/accessors';

	export default {
		name: "PublishedField",
		mixins: [Accessors],
		methods: {
			validate() {
				// make the parent validate the whole form to control
				// display of the submit button
				this.$bus.$emit("checkMetadata");
			}
		}
	};
</script>