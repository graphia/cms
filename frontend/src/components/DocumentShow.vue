<template>
	<div id="document_show">
		<div v-html="document.html"/>
	</div>
</template>

<script lang="babel">
	export default {
		created() {
			// populate $store.state.documents with docs from api

			let directory = this.$route.params.directory;
			let filename = this.$route.params.filename;
			
			console.debug(`retrieving document ${filename} from ${directory}`);
			
			this.$store.dispatch("getDocument", {directory, filename});
		},
		computed: {
			document() {
				return this.$store.state.activeDocument;
			}
		}
	}
</script>