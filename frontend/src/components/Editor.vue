<template>
	<div id="editor" class="form-group">
		<label for="markdown" class="sr-only">Document Contents</label>
		<textarea name="markdown" class="form-control" rows="40" v-model="document.markdown"/>
	</div>
</template>

<script lang="babel">
	export default {
		name: "Editor",
		computed: {
			document() {
				return this.$store.state.activeDocument;
			}
		}
	}
</script>