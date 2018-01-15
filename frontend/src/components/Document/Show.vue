<template>
	<div v-title="document.title">
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

							<dt>Date</dt>
							<dd>{{ document.date }}</dd>

							<dt>Tags</dt>
							<dd>
								<span v-for="(tag, i) in document.tags" class="tag badge badge-primary" :key="i">
									{{ tag }}
								</span>
							</dd>

							<dt>Draft</dt>
							<dd>{{ this.draftDescription() }}</dd>
						</dl>

						<div class="btn-toolbar" role="toolbar">

							<router-link class="btn btn-primary mr-2" :to="{name: 'document_edit', params: this.navigationParams}">
								Edit
							</router-link>

							<Translation v-if="$store.state.translationEnabled"/>

							<router-link class="btn btn-info mr-2" :to="{name: 'document_history', params: this.navigationParams}">
								History
							</router-link>

							<DocumentDelete/>

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
	import Translation from './Translation';
	import DocumentDelete from './Buttons/Delete';
	import CMSBreadcrumb from '../../javascripts/models/breadcrumb.js';
	import Accessors from '../Mixins/accessors';

	import checkResponse from "../../javascripts/response.js";

	export default {
		name: "DocumentShow",
		created() {
			// populate $store.state.documents with docs from api

			this.getDocument();

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

				let directory_title, doc;

				// if we have it, use the metadata provided directory and
				// document title
				if (this.document.directory_info) {
					directory_title = this.document.directory_info.title;
					doc = this.document.title;
				};

				return [
					new CMSBreadcrumb(
						directory_title || this.directory,
						"document_index",
						{directory: this.directory}
					),
					new CMSBreadcrumb(
						doc || this.params.document,
						"document_show",
						{directory: this.directory, document: doc}
					)
				];
			},
			navigationParams() {
				let p = {
					directory: this.params.directory,
					document: this.params.document
				};

				if (this.document.isTranslation()) {
					p.language_code = this.document.language;
				};

				return p;
			}
		},
		methods: {
			async getDocument() {

				let filename = "index.md";

				if (this.params.language_code) {
					filename = `index.${this.params.language_code}.md`;
				};

				let document = this.params.document;
				let directory = this.params.directory;

				this.$store.dispatch("getDocument", {directory, document, filename});
			},
			draftDescription() {
				if (this.document.draft) {
					return "Yes";
				};
				return "No";
			}
		},
		mixins: [Accessors],
		components: {
			Breadcrumbs,
			Translation,
			DocumentDelete
		}
	}
</script>