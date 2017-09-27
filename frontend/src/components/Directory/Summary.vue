<template>
	<div>

		<div class="jumbotron" v-if="numberOfDirectories == 0">
			{{ directories.length }}
			<p>No directories</p>

			<DirectoryNewButton/>

		</div>

		<!-- listing directories -->
		<div class="row" v-else-if="numberOfDirectories > 0">

			<div class="col-lg-4 mt-4" v-for="(directory, i) in directories" :key="i">

				<div class="card" :class="directory.path" >

					<h4 class="card-header">
						<router-link :to="{name: 'document_index', params: {directory: directory.path}}">
							{{ (directory.info.title || directory.path) | capitalize }}
						</router-link>
					</h4>

					<div class="card-body">
						{{ directory.info.description }}
					</div>

					<!-- listing documents inside a directory -->
					<div class="list-group list-group-flush" v-if="directory.contents.length > 0">

						<router-link
							v-for="(document, j) in directory.contents"
							:key="j"
							:to="{name: 'document_show', params: {directory: directory.path, filename: document.filename}}"
							:data-filename="document.filename"
							class="list-group-item list-group-item-action"
						>

							<h5>{{ document.frontmatter.title || document.filename }}</h5>

							<p class="text-muted">{{ document.frontmatter.synopsis || "No synopsis has been added" }}</p>

						</router-link>

						<div class="card-body">
							<router-link class="btn btn-sm btn-primary" :to="{name: 'document_new', params: {directory: directory.path}}">
								Create a document
							</router-link>
						</div>


					</div>
					<!-- /listing documents inside a directory -->

					<div class="card-body" v-else-if="directory.contents.length == 0">

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
	import CMSDirectory from '../../javascripts/models/directory.js';
	import DirectoryNewButton from './NewButton';

	export default {
		name: "DirectorySummary",
		data() {
			return {
				directories: {}
			}
		},
		created() {
			this.fetchDirectorySummary();
		},
		computed: {
			numberOfDirectories() {
				let count = Object.keys(this.directories).length;
				console.debug("directory count", count);
				return count;
			}
		},
		methods: {
			async fetchDirectorySummary() {

				let path = `${config.api}/summary`

				console.log("fetching directories", path);

				try {
					let response = await fetch(path, {
						mode: "cors",
						method: "GET",
						headers: this.$store.state.auth.authHeader()
					});

					if (!checkResponse(response)) {
						console.error(response);
						return;
					}

					let json = await response.json();
					console.log("got json", json)

					// TODO map the directories into CMSFile objects
					this.directories = json;
					console.log("directories", json);

				}
				catch(error) {
					console.error(error);
				}
			}
		},
		components: {
			DirectoryNewButton
		}
	}
</script>