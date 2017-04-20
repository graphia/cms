<template>
	<div id="document-index">
		<h2>Documents</h2>
		<ul>
			<li v-for="document in documents">
				<router-link :to="{name: 'document_show', params: {filename: document.filename}}">
					{{ document.absolutePath }}
				</router-link>
			</li>
		</ul>

		<router-link :to="{name: 'document_new', params: {directory: this.$route.params.directory}}" class="btn btn-primary">
			New
		</router-link>
	</div>
</template>

<script lang="babel">
	export default {
		name: "DocumentIndex",
		created() {
			// populate $store.state.documents with docs from api

			console.debug("retrieving all files from", this.directory);

			this.$store.dispatch("getDocumentsInDirectory", this.directory);
		},
		computed: {
			documents() {
				return this.$store.state.documents;
			},
			directory() {
				return this.$route.params.directory;
			}
		}
	}
</script>
