<template>
	<div class="row history">


		<div class="col-sm-12">

			<h1>History</h1>

			<ol class="commit-list">
				<li class="commit-list-item" v-for="item in history">
					<div class="card">

						<div class="card-header">
							{{ item.author.When | format_date }}
						</div>

						<div class="card-body">

							<p class="card-text">
								{{ item.message }}
							</p>

							<a class="card-link" :href="`mailto:${item.author.Email}`">{{ item.author.Name }}</a>
							<router-link class="card-link" :to="{name: 'commit', params: {hash: item.id}}">View entire commit</router-link>

						</div>
					</div>

				</li>
			</ol>

		</div>
	</div>
</template>

<script lang="babel">

	import checkResponse from '../../javascripts/response.js';

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

			this.history = await response.json()

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