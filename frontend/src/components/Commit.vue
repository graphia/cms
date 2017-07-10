<template>
	<div class="row">
		<div class="col-md-12 commit">
			<h1>Commit</h1>

			<dl class="row">

				<dt class="col-sm-3">Commit ID</dt>
				<dd class="col-sm-9">{{ this.commit.hash }}</dd>

				<dt class="col-sm-3">Author</dt>
				<dd class="col-sm-9"><a :href="`mailto:${this.committerEmailAddress}`">{{ this.committerName }}</a></dd>

				<dt class="col-sm-3">Message</dt>
				<dd class="col-sm-9">{{ this.commit.message }}</a></dd>

				<dt class="col-sm-3">Time</dt>
				<dd class="col-sm-9">{{ this.commit.timestamp | format_date }}</a></dd>

			</dl>

			<ol class="files">
				<li v-for="(item, key, index) in files">
					<CommitFile :path='key' :files='item'/>
				</li>
			</ol>

		</div>
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
			},
			committerName() {
				try {
					return this.commit.author.Name;
				} catch(err) {
					return "None found";
				}
			},
			committerEmailAddress() {
				try {
					return this.commit.author.Email;
				} catch(err) {
					return "None found";
				}
			}
		},
		created() {
			console.debug("created");
			this.retrievePatch(this.hash);
		},
		components: {
			CommitFile
		},
		methods: {
			async retrievePatch() {

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

<style lang="scss">
	div.commit {
		max-width: 60em;

		ol.files {
			list-style: none;
			padding-left: 0em;

			li {
				margin-bottom: 1em;
			}
		}
	}
</style>