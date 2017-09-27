<template>
	<div class="form-group document-version">

		<label for="document-version" class="form-control-label">Version</label>
		<input
			id="tag-editor"
			type="text"
			name="version"
			class="form-control"
			v-model="document.version"
		/>

	</div>
</template>

<script lang="babel">
	export default {
		name: "VersionField",
		computed: {
			document() {
				return this.$store.state.activeDocument;
			}
		}
	}
</script>
