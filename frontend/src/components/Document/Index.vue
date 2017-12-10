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

				<div class="col-md-4" v-for="(d, base, i) in groupedTranslations" :key="i">

					<div class="card document-entry m-1" :data-filename="base">

						<h3 class="card-header">
							<router-link :to="{name: 'document_show', params: {filename: primary(d).filename}}">
								{{ primary(d).title || primary(d).filename }}
							</router-link>
						</h3>

						<div class="card-body">
							<p class="card-text">{{ primary(d).synopsis || description_placeholder }}</p>
						</div>

						<div class="card-footer" v-if="translationEnabled && d.length > 1">
							<ul class="list-inline">
								<li class="list-inline-item" v-for="(t, k) in translations(d)" :key="k" :data-lang="t.language.name">
									<router-link :to="{name: 'document_show', params: {filename: t.filename}}">
										{{ (t.language && t.language.flag) || "missing" }}
									</router-link>
								</li>
							</ul>
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
	import Accessors from '../Mixins/accessors';

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
				this.$store.dispatch("getDocumentsInDirectory", directory);
			},

			// return the 'primary' (first) copy of a file,
			// usually in the default language
			primary(files) {
				return files[0];
			},

			translations(files) {
				return files
					.filter((file) => { return file.translation })
			}


		},
		computed: {
			title() {
				return this.$store.activeDirectory.title;
			},
			breadcrumbs() {
				return [
					new CMSBreadcrumb(
						this.activeDirectory.title || this.directory,
						"document_index",
						{directory: this.directory}
					)
				];
			},
			translationEnabled() {
				return this.$store.state.translationEnabled;
			},
			groupedTranslations() {

				// FIXME finding translations using the number of dots is potentially
				// a bit fragile, should use a regexp to check for filenames ending
				// in ".xx.md"
				return this
					.documents
					.sort((a,b) => {
						// default language files first
						return (a.filename.split(".").length - b.filename.split(".").length)
					})
					.reduce((summary, doc) => {
						// use the file's basename to group translations
						let base = doc.filename.split(".")[0]

						summary[base] ? summary[base].push(doc) : summary[base] = [doc];

						return summary;
					}, {});
			}
		},
		mixins: [Accessors],
		components: {
			Breadcrumbs,
			Error
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