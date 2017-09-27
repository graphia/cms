<template>
	<div id="document-index">

		<div v-if="this.documents && this.documents.length > 0">

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
								{{ document.title || document.filename }}
							</router-link>
						</h3>

						<div class="card-body">
							<p class="card-text">{{ document.synopsis || description_placeholder }}</p>
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

		<div v-else-if="this.documents && this.documents.length === 0">

			<div class="col-12">

				<div class="alert alert-warning">

					<h3>There's nothing here <em>yet!</em></h3>

					<p>
						This directory is empty. Don't worry, you can add the first document by clicking the button below.
					</p>

					<router-link :to="{name: 'document_new', params: {directory: this.$route.params.directory}}" class="btn btn-primary">
						Create a new document
					</router-link>
				</div>
			</div>
		</div>

		<div v-else>
			<Error :code="404"/>
		</div>
	</div>
</template>

<script lang="babel">

	import Breadcrumbs from '../Utilities/Breadcrumbs';
	import Error from '../Errors/Error';
	import CMSBreadcrumb from '../../javascripts/models/breadcrumb.js';

	export default {
		name: "DocumentIndex",
		created() {
			// populate $store.state.documents with docs from api
			this.setup(this.directory);
		},
		watch: {
			// if we navigate from one dir index to another, reload the
			// contents
			"$route"(to, from) {
				this.setup(this.directory);
			}
		},
		data() {
			return {
				description_placeholder: "No description has been added"
			};
		},
		methods: {
			async setup(directory) {
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
						this.activeDirectory.title || this.directory,
						"document_index",
						{directory: this.directory}
					)
				];
			}
		},
		components: {
			Breadcrumbs,
			Error
		}
	}
</script>

<style lang="scss">
</style>