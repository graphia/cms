<template>
	<div>

		<div class="col col-md-6" v-if="loading">
			<div class="alert alert-warning mx-auto">
				Loading...
			</div>
		</div>

		<div class="col col-md-6" v-else-if="numberOfDirectories == 0">
			<div class="jumbotron">
				{{ directories.length }}
				<p>No directories</p>

				<DirectoryNewButton/>
			</div>
		</div>

		<!-- listing directories -->
		<div class="row" v-else-if="numberOfDirectories > 0">

			<div class="col-lg-4 mt-4" v-for="(directory, i) in directoriesWithGroupedTranslations" :key="i">

				<div class="card" :class="directory.path" :data-directory="directory.path">

					<h4 class="card-header">
						<router-link :to="{name: 'document_index', params: {directory: directory.path}}">
							{{ (directory.info.title || directory.path) | capitalize }}
						</router-link>
					</h4>

					<div class="card-body">
						{{ directory.info.description }}
					</div>

					<!-- listing documents inside a directory -->
					<div class="list-group list-group-flush" v-if="Object.keys(directory.contents).length > 0">

						<router-link
							v-for="(documents, base, j) in directory.contents"
							:key="j"
							:to="{name: 'document_show', params: {directory: directory.path, document: primary(documents).document, filename: primary(documents).filename}}"
							:data-filename="base"
							class="list-group-item list-group-item-action"
						>

							<h5>{{ primary(documents).title || primary(documents).filename }}</h5>

							<p class="text-muted">{{ primary(documents).synopsis || "No synopsis has been added" }}</p>

							<div class="translations" v-if="translationEnabled">

								<ul class="translations-list list-inline">
									<li class="list-inline-item" v-for="(t, k) in translations(documents)" :key="k" :data-lang="t.language.name">
										<router-link :to="{name: 'document_show', params: {directory: directory.path, filename: t.filename}}">
											{{ (t.language && t.language.flag) || "missing" }}
										</router-link>
									</li>
								</ul>

							</div>

						</router-link>

						<div class="card-body">
							<router-link class="btn btn-sm btn-primary" :to="{name: 'document_new', params: {directory: directory.path}}">
								Create a document
							</router-link>
						</div>


					</div>
					<!-- /listing documents inside a directory -->

					<div class="card-body" v-else>

						<div class="alert alert-info">
							There's nothing here yet
						</div>

						<router-link class="btn btn-sm btn-primary" :to="{name: 'document_new', params: {directory: directory.path}}">
							Create a document
						</router-link>

					</div>

				</div>
			</div>

			<div class="col-lg-4 mt-4">
				<DirectoryNewButton/>
			</div>

		</div>
		<!-- /listing directories -->

	</div>
</template>

<script lang="babel">

	import 'babel-runtime/core-js/object/keys';

	import checkResponse from '../../javascripts/response.js';
	import config from '../../javascripts/config.js';
	import CMSFile from '../../javascripts/models/file.js';
	import CMSDirectory from '../../javascripts/models/directory.js';

	import DirectoryNewButton from './NewButton';

	export default {
		name: "DirectorySummary",
		data() {
			return {
				directories: {},
				loading: true
			}
		},
		created() {
			this.fetchDirectorySummary();
		},
		computed: {
			numberOfDirectories() {
				let count = Object.keys(this.directories).length;
				return count;
			},
			translationEnabled() {
				return this.$store.state.translationEnabled;
			},
			directoriesWithGroupedTranslations() {

				return this.directories
					.map((dir) => {

						let groupedTranslations = dir
							.contents
							.sort((a,b) => {
								// default language files first
								return (a.filename.split(".").length - b.filename.split(".").length)
							})
							.reduce((summary, doc) => {

								// use the file's basename to group translations
								//let base = doc.filename.split(".")[0]
								let base = doc.document;
								let parsedDoc = new CMSFile(doc);

								summary[base] ? summary[base].push(parsedDoc) : summary[base] = [parsedDoc];

								return summary;
							}, {});

						return {
							path: dir.path,
							info: dir.info,
							contents: groupedTranslations
						};

					});

			}
		},
		methods: {
			async fetchDirectorySummary() {

				let path = `${config.api}/summary`

				try {
					let response = await fetch(path, {
						mode: "cors",
						method: "GET",
						headers: this.$store.state.auth.authHeader()
					});

					if (!checkResponse(response.status)) {
						console.error(response);
						return;
					}

					let json = await response.json();

					this.loading = false;

					// TODO map the directories into CMSFile objects
					this.directories = json;

				}
				catch(error) {
					console.error(error);
				};

			},
			primary(files) {
				return files[0];
			},

			translations(files) {
				return files
					.filter((file) => { return file.translation })
			}
		},
		components: {
			DirectoryNewButton
		}
	}
</script>

<style lang="scss">
	.translations > .translations-list {
		margin-bottom: 0px;
	}
</style>