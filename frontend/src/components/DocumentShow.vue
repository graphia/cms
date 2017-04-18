<template>
	<section class="row">

		<article id="document_show" class="col-md-9">
			<div v-html="document.html"/>
		</article>

		<div class="col-md-3">
			<div class="btn-toolbar">
				<router-link :to="{name: 'document_edit', params: {directory: 'documents', filename: document.filename}}" class="btn btn-success">
					Edit
				</router-link>
			</div>
		</div>
	</section>
</template>

<script lang="babel">
	export default {
		name: "DocumentShow",
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