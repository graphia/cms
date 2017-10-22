<template>
	<div>
		<Breadcrumbs :levels="breadcrumbs"/>

		<section class="row">

			<article id="document-content" class="col-md-8">
				<div class="content" v-html="relativeHTML"/>
			</article>

			<aside class="col-md-4">
				<div class="card document-metadata">
					<div class="card-body">
						<dl>

							<dt>Title</dt>
							<dd>{{ document.title  }}</dd>


							<dt>Description</dt>
							<dd>{{ document.synopsis }}</dd>

							<dt>Version</dt>
							<dd>{{ document.version }}</dd>

							<dt>Author</dt>
							<dd>{{ document.author }}</dd>

							<dt>Tags</dt>
							<dd>
								<span v-for="tag in document.tags" class="tag badge badge-primary">
									{{ tag }}
								</span>
							</dd>
						</dl>

						<div class="btn-toolbar" role="toolbar">
							<router-link class="btn btn-success mr-2" :to="{name: 'document_edit', params: {directory: this.directory, filename: this.filename}}">
								Edit
							</router-link>

							<router-link class="btn btn-info mr-2" :to="{name: 'document_history', params: {directory: this.directory, filename: this.filename}}">
								History
							</router-link>

							<button type="button" @click="destroy" class="btn btn-danger mr-2">
								Delete
							</button>

						</div>
					</div>
				</div>
			</aside>
		</section>
	</div>
</template>

<style scoped lang="scss">
	aside {
		margin: 2em 0em;
	}

	.document-metadata {
		span.tag {
			margin-right: 0.6em;
		}
	}
</style>

<script lang="babel">

	import Breadcrumbs from '../Utilities/Breadcrumbs';
	import CMSBreadcrumb from '../../javascripts/models/breadcrumb.js';
	import Accessors from '../Mixins/accessors';

	export default {
		name: "DocumentShow",
		created() {
			// populate $store.state.documents with docs from api

			let directory = this.directory;
			let filename = this.filename;

			console.debug(`retrieving document ${filename} from ${directory}`);

			this.$store.dispatch("getDocument", {directory, filename});

			// create a commit to be populated/used if delete is clicked
			this.$store.dispatch("initializeCommit")

		},
		computed: {

			// Amend any relative links or images to point at the
			// correct resource
			relativeHTML() {

				let attachmentsDir = this.document.slug;
				let html = $.parseHTML(this.document.html);

				$(html)
					.find('img')
					.each(function(_, image) {
						if ($(image)
							.attr('src')
							.startsWith("images")) {
								let src = $(image).attr('src');
								$(image).attr('src', [attachmentsDir, src].join("/"));
						};
					});

				return html
						.map((e) => {return e.outerHTML})
						.join("");


			},
			breadcrumbs() {

				let directory_title, filename;

				// if we have it, use the metadata provided directory and
				// document title
				if (this.document.directory_info) {
					directory_title = this.document.directory_info.title;
					filename = this.document.title;
				};

				return [
					new CMSBreadcrumb(
						directory_title || this.directory,
						"document_index",
						{directory: this.directory}
					),
					new CMSBreadcrumb(
						filename || this.filename,
						"document_show",
						{directory: this.directory, document: (filename || this.filename)}
					)
				];
			}
		},
		methods: {
			async destroy(event) {
				event.preventDefault();
				console.debug("delete clicked!");

				let file = this.document;

				let response = await this.document.destroy(this.commit);

				console.debug("File deleted, redirecting to document index");
				this.redirectToDirectoryIndex(this.directory);

			},
			redirectToDirectoryIndex(directory) {
				this.$router.push({
					name: 'document_index',
					params:{directory}
				});
			}
		},
		mixins: [Accessors],
		components: {
			Breadcrumbs
		}
	}
</script>