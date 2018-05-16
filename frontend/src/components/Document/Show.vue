<template>
	<div v-title="document.title || 'Not found'">
		<Breadcrumbs :levels="breadcrumbs"/>

		<section class="row" v-if="this.document && !this.document.initializing">

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

							<dt>Translations</dt>
							<dd class="translations">
								<ul class="list-inline">
									<li class="list-inline-item" v-for="(translation, i) in translations" :key="i" :data-lang="translation.code">
										<router-link :to="{name: 'document_show', params: translation.params}">
											{{ translation.flag || translation.code }}
										</router-link>
									</li>
								</ul>
							</dd>
						</dl>

						<div class="btn-toolbar" role="toolbar">

							<router-link class="btn btn-primary my-2 mx-1" :to="{name: 'document_edit', params: this.navigationParams}">
								Edit
							</router-link>

							<Translation v-if="$store.state.server.translationInfo.translationEnabled"/>

							<router-link class="btn btn-info my-2 mx-1" :to="{name: 'document_history', params: this.navigationParams}">
								History
							</router-link>

							<DocumentDelete/>

						</div>
					</div>


				</div>
			</aside>
		</section>
		<div v-else>
			<Error :code="404"/>
		</div>
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
	import Error from '../Errors/Error';
	import DocumentDelete from './Buttons/Delete';
	import CMSBreadcrumb from '../../javascripts/models/breadcrumb.js';
	import Accessors from '../Mixins/accessors';

	import config from '../../javascripts/config.js';
	import checkResponse from "../../javascripts/response.js";
	import filenameFromLanguageCode from '../../javascripts/utilities/filename-from-language-code.js';

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

				let html = $.parseHTML(this.document.html);
				let dir = this.params.directory;
				let doc = this.params.document;

				$(html)
					.find('img')
					.each(function(_, image) {
						if ($(image)
							.attr('src')
							.startsWith("images")) {
								let src = $(image).attr('src');

								// use the absolute path so we don't need to worry about
								// translations which now have the path in the /cms/dir/doc/en style
								$(image).attr('src', [config.cms, dir, doc, src].join("/"));
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
						directory_title || this.params.directory,
						"directory_index",
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
			},

			translations() {

				const ld = this.$store.state.server.translationInfo.languages;
				const dl = this.$store.state.server.translationInfo.defaultLanguage;

				let tr = this.document.translations.map((code) => {
					let match = ld.find(x => x.code === code);

					if (!match) {
						console.warn("no language configured for code", code);
						return {code: code, flag: "", name: code};
					};

					return match;

				});

				// add in route params so we don't clutter the view
				return tr.map((t) => {
					// {directory: document.path, filename: document.filename, document: document.document, language_code: translation.code}
					t.params = {
						directory: this.document.path,
						filename: this.document.filename,
						document: this.document.document,
						language_code: (dl === t.code ? null : t.code)
					}
					return t;
				})
			}
		},
		methods: {
			async getDocument() {
				const filename = filenameFromLanguageCode(this.params.language_code)

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
			DocumentDelete,
			Error
		}
	}
</script>
