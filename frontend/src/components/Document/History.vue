<template>
	<div class="row history">

		<div class="col-sm-12">

			<Breadcrumbs :levels="breadcrumbs"/>

			<h1>History</h1>

			<ol class="commit-list">
				<li class="commit-list-item" v-for="(item, i) in history" :key="i">
					<div class="card" :class="`commit-${item.id}`">

						<div class="card-header">
							{{ item.author.When | format_date }}
						</div>

						<div class="card-body">

							<p class="card-text">
								{{ item.message }}
							</p>

							<div class="btn-toolbar">
								<a class="card-link btn btn-secondary" :href="`mailto:${item.author.Email}`">{{ item.author.Name }}</a>

								<button type="button"
										class="btn btn-info"
										data-toggle="collapse"
										:data-target="`#diff-${item.id}`"
								>
									Show changes
								</button>

								<router-link class="card-link btn btn-info" :to="{name: 'commit', params: {hash: item.id}}">View entire commit</router-link>
							</div>


							<div class="collapse multi-collapse" :id="`diff-${item.id}`">
								<pre v-html="item.patch.diff()"/>
							</div>

						</div>

					</div>

				</li>
			</ol>

		</div>
	</div>
</template>

<script lang="babel">

	import Breadcrumbs from '../Utilities/Breadcrumbs';

	import checkResponse from '../../javascripts/response.js';
	import CMSBreadcrumb from '../../javascripts/models/breadcrumb.js';
	import CMSPatch from '../../javascripts/models/patch.js';


	import Diff from 'text-diff';

	export default {
		name: "DocumentHistory",

		data() {
			return {
				history: []
			};
		},

		computed: {
			directory() {
				return this.$route.params.directory;
			},
			filename() {
				return this.$route.params.filename;
			},
			document() {
				return this.$store.state.activeDocument;
			},
			breadcrumbs() {
				let directory_title, filename;

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
						{directory: this.document.path, document: this.document.filename}
					),
					new CMSBreadcrumb(
						"History",
						"document_history",
						{directory: this.directory, document: (filename || this.filename)}
					)
				];
			}
		},

		async created() {

			var directory = this.directory;
			var filename = this.filename;

			if (!this.$store.state.activeDocument.populated()) {
				await this.$store.dispatch("getDocument", {directory, filename});
			};

			let response = await this.$store.state.activeDocument.log();

			if (!checkResponse(response.status)) {
				throw("Could not retrieve history");
			}

			let json = await response.json();

			this.history = json.map((revision) => {
				revision.patch = new CMSPatch(revision.id, "", revision.old, revision.new);
				return revision;
			});

		},
		components: {
			Breadcrumbs
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