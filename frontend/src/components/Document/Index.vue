<template>
	<div id="document-index">

		<div class="row document-info">
			<div class="col-md-12">

				<!-- document index header -->
				<h2 v-if="activeDirectory.title">
					{{ activeDirectory.title }}
				</h2>
				<h2 v-else>
					{{ directory | capitalize }}
				</h2>
				<!-- /document index header -->

				<p>{{ activeDirectory.description }}</p>

			</div>
		</div>


		<div class="row document-list card-deck">

			<div class="card col col-md-6 col-lg-3 p-3 m-1" v-for="(document, i) in documents" :key="i">

					<div class="card-body">

						<h4 class="card-title">
							{{ document.title }}
						</h4>

						<p>{{ document.synopsis }}</p>

						<router-link :to="{name: 'document_show', params: {filename: document.filename}}">
							{{ document.absolutePath }}
						</router-link>
					</div>

			</div>

		</div>

		<div class="row document-buttons">

			<div class="col-12">
				<router-link :to="{name: 'document_new', params: {directory: this.$route.params.directory}}" class="btn btn-primary">
					New
				</router-link>
			</div>

		</div>
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
