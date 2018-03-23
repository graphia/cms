<template>
	<div class="mt-4">

		<div class="col col-md-8" v-if="loading">
			<div class="alert alert-warning mx-auto">
				Loading...
			</div>
		</div>

		<div class="col col-md-4" v-else-if="directories.length == 0">
			<div class="alert alert-info">
				Your repository is still empty! Once created, your documents
				will appear here!
			</div>
		</div>

		<div class="col col-md-12 directories" v-else>

			<!-- listing directories -->
				<div class="border rounded p-4 mb-4 bg-white directory" :data-directory="directory.path" :class="directory.path" v-for="(directory, i) in directories" :key="i">

					<div class="row">
						<h4 class="col col-md-6">
							<router-link :to="{name: 'directory_index', params: {directory: directory.path}}">
								{{ (directory.title || directory.path) | capitalize }}
							</router-link>
						</h4>

						<div class="col col-md-6 text-right">
							<DocumentNewButton :directoryPath="directory.path"/>
						</div>
					</div>

					<p class="mt-4 mb-4 directory-description">{{ directory.description }}</p>

					<IndexList :documents="directory.contents" :directoryPath="directory.path"/>

				</div>

			<!-- /listing directories -->
		</div>

	</div>
</template>

<script lang="babel">

	import 'babel-runtime/core-js/object/keys';

	import checkResponse from '../../javascripts/response.js';
	import config from '../../javascripts/config.js';
	import CMSFile from '../../javascripts/models/file.js';
	import CMSDirectory from '../../javascripts/models/directory.js';
	import IndexList from '../Directory/Index/List';
	import DocumentNewButton from '../Document/Buttons/New';

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
			this.$store.commit("refreshServerInfo");
		},
		computed: {
			numberOfDirectories() {
				let count = Object.keys(this.directories).length;
				return count;
			}
		},
		methods: {
			async fetchDirectorySummary() {

				let path = `${config.api}/summary`

				let response = await fetch(path, {
					method: "GET",
					headers: this.$store.state.auth.authHeader()
				});

				if (!checkResponse(response.status)) {
					console.error(response);
					return;
				}

				let json = await response.json();

				this.directories = json.map((dir) => {

					return new CMSDirectory(
						dir.path,
						dir.info.title,
						dir.info.description,
						dir.info.body,
						dir.info.html,
						dir.contents.map((file) => {
							return new CMSFile(file);
						})
					);
				});

				this.loading = false;

			}
		},
		components: {
			IndexList,
			DocumentNewButton
		}
	}
</script>

<style lang="scss">
	.translations > .translations-list {
		margin-bottom: 0px;
	}
</style>