<template>
	<div id="document-index">
		<h2>Documents</h2>

		<ul>
			<li v-for="document in documents">
				{{ document.absolutePath }}
			</li>
		</ul>
	</div>
</template>

<script lang="babel">
	export default {
		name: "DocumentIndex",
		created() {
			// populate $store.state.documents with docs from api
			return this.$store.dispatch("getDocumentsInDirectory", "documents");
		},

		computed: {
			documents() {
				console.log("accessing stored docs");
				return this.$store.state.documents;
			}
		}
	}
</script>
