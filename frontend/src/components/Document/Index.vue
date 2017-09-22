<template>
	<div id="document-index">

		<Breadcrumbs :levels="breadcrumbs"/>

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


		<div class="row document-list">

			<div class="col-md-4" v-for="(document, i) in documents" :key="i">

				<div class="card m-4" :data-filename="document.filename">

					<h3 class="card-header">
						<router-link :to="{name: 'document_show', params: {filename: document.filename}}">
							{{ document.title }}
						</router-link>
					</h3>

					<div class="card-body">
						<p class="card-text">{{ document.synopsis }}</p>
					</div>

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

	import Breadcrumbs from '../Utilities/Breadcrumbs';
	import CMSBreadcrumb from '../../javascripts/models/breadcrumb.js';

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
			title() {
				return this.$store.activeDirectory.title;
			},
			documents() {
				return this.$store.state.documents;
			},
			directory() {
				return this.$route.params.directory;
			},
			activeDirectory() {
				return this.$store.state.activeDirectory;
			},
			breadcrumbs() {
				return [
					new CMSBreadcrumb(
						this.activeDirectory.title,
						"document_index",
						{directory: "pokemon"}
					)
				];
			}
		},
		components: {
			Breadcrumbs
		}
	}
</script>

<style lang="scss">
</style>