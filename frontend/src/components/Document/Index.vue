<template>
	<div id="document-index">
		<h2 v-if="activeDirectory.title">
			{{ activeDirectory.title }}
		</h2>
		<h2 v-else>
			{{ directory | capitalize }}
		</h2>

		<p>{{ activeDirectory.description }}</p>

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
			this.fetchDocuments(this.directory);
		},
		watch: {
			// if we navigate from one dir index to another, reload the
			// contents
			"$route"(to, from) {
				console.log("changed");
				this.fetchDocuments(this.directory);
			}
		},
		methods: {
			fetchDocuments(directory) {
				console.debug("retrieving all files from", directory);
				this.$store.dispatch("getDocumentsInDirectory", directory);
			}
		},
		computed: {
			documents() {
				return this.$store.state.documents;
			},
			directory() {
				return this.$route.params.directory;
			},
			activeDirectory() {
				return this.$store.state.activeDirectory;
			}
		}
	}
</script>
