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
	</div>
</template>

<script lang="babel">
	export default {
		name: "DocumentIndex",
		created() {
			// populate $store.state.documents with docs from api
			let directory = this.$route.params.directory;

			console.debug("retrieving all files from", directory);

			this.$store.dispatch("getDocumentsInDirectory", directory);
		},

		computed: {
			documents() {
				return this.$store.state.documents;
			}
		}
	}
</script>
