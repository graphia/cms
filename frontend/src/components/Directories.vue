<template>
	<div>
		<h1>New Directory</h1>


		<div v-if="numberOfDirectories == 0">
			{{ directories.length }}
			<p>No directories</p>

		</div>

		<div v-else-if="numberOfDirectories > 0">



			<div class="card" v-for="(contents, directory) in directories">

				<div class="card-header">
					<h3 class="card-title">{{ directory | capitalize }}</h3>
				</div>

				<div class="list-group list-group-flush">

					<a href="#" class="list-group-item list-group-item-action" v-for="document in contents">
						{{ document.frontmatter.title }}
					</a>

				</div>

			</div>
		</div>

		<div class="new-directory">

			<form>

				<div class="input-group">
					<input class="form-control" placeholder="stories"/>

					<span class="input-group-btn">
						<input type="submit" class="form-control btn btn-success" @click="createDirectory" value="Create Directory"/>
					</span>

				</div>
			</form>

		</div>

	</div>
</template>

<script lang="babel">

	import checkResponse from '../javascripts/response.js';
	import config from '../javascripts/config.js';
	import _ from 'babel-runtime/core-js/object/keys';

	export default {
		name: "Directories",
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
			createDirectory(event) {
				event.preventDefault();

				console.debug("clicked");
			},
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
		}
	}
</script>