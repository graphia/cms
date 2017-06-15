<template>
	<div>
		<h1>Commit</h1>

		<dl class="row">

			<dt class="col-sm-3">Commit ID</dt>
			<dd class="col-sm-9">{{ this.commit.hash }}</dd>

			<dt class="col-sm-3">Author</dt>
			<dd class="col-sm-9"><a :href="`mailto:${this.commit.author.Email}`">{{ this.commit.author.Name }}</a></dd>

			<dt class="col-sm-3">Message</dt>
			<dd class="col-sm-9">{{ this.commit.message }}</a></dd>

		</dl>

		<ol>
			<li v-for="(item, key, index) in files">
				<CommitFile :path='key' :files='item'/>
			</li>
		</ol>


	</div>

</template>

<script lang="babel">
	import config from '../javascripts/config.js';
	import store from '../javascripts/store.js';
	import checkResponse from '../javascripts/response.js';
	import CommitFile from './Commit/File';

	export default {
		name: "Commit",
		data() {
			return {
				commit: {}
			}
		},
		computed: {
			hash() {
				return this.$route.params.hash;
			},
			files() {
				return this.commit.files;
			}
		},
		async created() {
			console.debug("created");
			await this.retrievePatch(this.hash);
			this.setupDiff();
		},
		components: {
			CommitFile
		},
		methods: {
			async retrievePatch() {
				console.debug("got here");

				let path = `${config.api}/commits/${this.hash}`;

				let response = await fetch(path, {mode: "cors", headers: store.state.auth.authHeader()});

				if (!checkResponse(response.status)) {
					throw("Could not retrieve changeset");
				}

				let json = await response.json()

				this.commit = json;


			}
		}
	};
</script>