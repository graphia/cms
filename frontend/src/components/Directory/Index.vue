<template>
	<div id="document-index" v-title="title">

		<div v-if="this.documents && this.documents.length > 0">

			<Breadcrumbs :levels="breadcrumbs"/>

			<div class="row directory-info">
				<div class="col-md-6">

					<!-- document index header -->
					<h2 v-if="activeDirectory.title">
						{{ activeDirectory.title }}
					</h2>
					<h2 v-else>
						{{ directory | capitalize }}
					</h2>
					<!-- /document index header -->

					<p>{{ activeDirectory.description }}</p>

					<div class="directory-info-text">
						<h3>Extra Information</h3>
						<div v-html="activeDirectory.html"/>
					</div>

				</div>

				<div id="directory-toolbar" class="col col-md-6 text-right">
					<DocumentNewButton :directoryPath="directory"/>
					<router-link :to="{name: 'directory_edit', params: {directory: this.$route.params.directory}}" class="btn btn-sm btn-primary">
						Edit directory
					</router-link>
					<DirectoryDeleteButton/>
				</div>
			</div>

			<IndexList :documents="documents" :directoryPath="directory"/>

		</div>

		<div v-else-if="this.documents && this.documents.length === 0">

			<Breadcrumbs :levels="breadcrumbs"/>

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
			<Breadcrumbs :levels="breadcrumbs"/>

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