<template>
	<div id="document-index" v-title="title">

		<Breadcrumbs :levels="breadcrumbs"/>

		<div class="rounded p-4 mb-4 bg-white" v-if="this.documents && this.documents.length > 0">

			<div class="row directory-info">
				<div class="col-md-12">

					<!-- document index header -->
					<h1 v-if="activeDirectory.title">
						{{ activeDirectory.title }}
					</h1>
					<h1 v-else>
						{{ directory | capitalize }}
					</h1>
					<!-- /document index header -->

					<blockquote class="blockquote directory-description">{{ activeDirectory.description }}</blockquote>

					<div class="directory-info-text" v-if="activeDirectory.html.length > 0">
						<div v-html="activeDirectory.html"/>
					</div>

				</div>


				<div id="directory-toolbar" class="col-md-12">
					<div class="mx-1 my-2">
						<DocumentNewButton :directoryPath="directory"/>
						<router-link :to="{name: 'directory_edit', params: {directory: this.$route.params.directory}}" class="btn btn-sm btn-secondary">
							Edit directory
						</router-link>
						<DirectoryDeleteButton/>
					</div>
				</div>

			</div>

			<IndexList :documents="documents" :directoryPath="directory"/>

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

	import CMSBreadcrumb from '../../javascripts/models/breadcrumb.js';

	import IndexList from './Index/List';
	import Breadcrumbs from '../Utilities/Breadcrumbs';
	import Error from '../Errors/Error';
	import Accessors from '../Mixins/accessors';
	import DocumentNewButton from '../Document/Buttons/New';
	import DirectoryDeleteButton from './Buttons/Delete';

	export default {
		name: "DirectoryIndex",
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
				this.$store.dispatch("getDocumentsInDirectory", directory);
			},

			// return the 'primary' (first) copy of a file,
			// usually in the default language


		},
		computed: {
			title() {
				return (this.activeDirectory && this.activeDirectory.title) || "Listing documents";
			},
			breadcrumbs() {
				return [
					new CMSBreadcrumb(
						this.activeDirectory.title || this.directory,
						"directory_index",
						{directory: this.directory}
					)
				];
			},


		},
		mixins: [Accessors],
		components: {
			Breadcrumbs,
			Error,
			IndexList,
			DocumentNewButton,
			DirectoryDeleteButton
		}
	}
</script>

<style lang="scss">
	.card.document-entry .card-footer {
		padding: 0.2rem 1.25rem;
		ul {
			margin-bottom: 0;
		}
	}
</style>