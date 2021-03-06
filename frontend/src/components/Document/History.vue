<template>
	<div class="row history" v-title="title">

		<div class="col-sm-12">

			<Breadcrumbs :levels="breadcrumbs"/>

			<h1>History</h1>

			<ol class="commit-list">
				<li class="commit-list-item" v-for="(item, i) in history" :key="i">
					<div class="card history" :class="`commit-${item.id}`">

						<div class="card-header">
							{{ item.author.When | format_date }}
						</div>

						<div class="card-body">

							<div class="card-text commit-message">

								<div class="mb-4" v-for="(line, j) in formatMessage(item.message)" :key="j">
									<h4 v-if="j == 0">
										{{ line }}
									</h4>
									<p v-else>
										{{ line }}
									</p>
								</div>

							</div>

							<div class="btn-toolbar">

								<a class="card-link btn btn-secondary" :href="`mailto:${item.author.Email}`">{{ item.author.Name }}</a>

								<router-link class="card-link btn btn-info ml-1" :to="{name: 'commit', params: {hash: item.id}}">View entire commit</router-link>

								<button type="button"
										class="btn btn-info ml-1"
										data-toggle="collapse"
										:data-target="`#diff-${item.id}`"
								>
									Show changes

									<octicon :icon-name="'chevron-down'"/>
								</button>

							</div>

							<Diff :patch="item.patch" :collapsible="true" :hash="item.id"/>

						</div>

					</div>

				</li>
			</ol>

		</div>
	</div>
</template>

<script lang="babel">

	import Breadcrumbs from '../Utilities/Breadcrumbs';
	import Diff from '../Utilities/Diff';
	import Accessors from '../Mixins/accessors';

	import checkResponse from '../../javascripts/response.js';
	import CMSBreadcrumb from '../../javascripts/models/breadcrumb.js';
	import CMSPatch from '../../javascripts/models/patch.js';
	import filenameFromLanguageCode from '../../javascripts/utilities/filename-from-language-code.js';

	export default {
		name: "DocumentHistory",

		data() {
			return {
				history: []
			};
		},

		methods: {
			formatMessage(m) {
				return m.split(/(\r\n|\n|\r)/gm);
			}
		},

		computed: {
			title() {
				return `${this.document.title}: History`;
			},
			breadcrumbs() {
				let directory_title, doc;

				if (this.document.directory_info) {
					directory_title = this.document.directory_info.title;
					doc = this.document.title;
				};

				return [

					new CMSBreadcrumb(
						directory_title || this.directory,
						"directory_index",
						{directory: this.directory}
					),
					new CMSBreadcrumb(
						doc || this.params.document,
						"document_show",
						{directory: this.params.directory, document: this.params.document}
					),
					new CMSBreadcrumb(
						"History",
						"document_history",
						{directory: this.params.directory, document: doc}
					)
				];
			}
		},

		async created() {

			const filename = filenameFromLanguageCode(this.params.language_code)

			var directory = this.params.directory;
			var document = this.params.document;

			try {

				if (!this.$store.state.activeDocument.populated()) {
					await this.$store.dispatch("getDocument", {directory, document, filename});
				};

				let response = await this.$store.state.activeDocument.log();

				if (!checkResponse(response.status)) {
					throw(`request failed ${response}`);
				};

				let json = await response.json();

				// create a CMSPatch (which performs the actual diff) and add it to the object
				this.history = json.map((revision) => {
					revision.patch = new CMSPatch(revision.id, "", revision.old, revision.new);
					return revision;
				});

			} catch(err) {
				console.error("Failed to get file history", err);
			};

		},
		mixins: [Accessors],
		components: {
			Breadcrumbs,
			Diff
		}
	};
</script>

<style lang="scss">

	div.history {
		max-width: 60em;

		ol.commit-list {
			list-style: none;
			padding-left: 0em;

			li {
				margin-bottom: 1em;
			}
		}
	}
</style>